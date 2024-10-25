package middleware

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/zstd"

	"compress/gzip"
	"github.com/klauspost/compress/flate"
	"github.com/klauspost/compress/snappy"
	"github.com/klauspost/compress/zlib"
	"github.com/pierrec/lz4"
	"github.com/ulikunitz/xz"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// CreateProxyHandler 创建代理处理函数
func CreateProxyHandler(
	upstreamSite string,
	processResponseBody func([]byte) ([]byte, error),
) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 构建上游服务器的 URL
		upstreamURL := upstreamSite + c.Request.URL.Path

		// 读取客户端请求的 Body
		reqBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("Failed to read request body, %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		err = c.Request.Body.Close()
		if err != nil {
			fmt.Printf("Failed to close request body, %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 创建新的请求发送到上游服务器
		req, err := http.NewRequest(c.Request.Method, upstreamURL, bytes.NewReader(reqBody))
		if err != nil {
			fmt.Printf("Failed to create request, %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 复制请求头
		req.Header = c.Request.Header.Clone()

		// 创建 HTTP 客户端并发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send request, %v", err)
			c.AbortWithStatus(http.StatusBadGateway)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Printf("Failed to close response body, %v", err)
				return
			}
		}(resp.Body)

		// 读取并解压缩响应体
		respBody, contentEncoding, err := readResponseBody(resp)
		if err != nil {
			fmt.Printf("Failed to read response body, %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 处理响应体
		modifiedBody, err := processResponseBody(respBody)
		if err != nil {
			fmt.Printf("Failed to process response body, %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 将响应返回给客户端
		err = writeResponse(c, resp, modifiedBody, contentEncoding)
		if err != nil {
			fmt.Printf("Failed to write response, %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

// 读取响应体并解压缩
func readResponseBody(resp *http.Response) ([]byte, string, error) {
	contentEncoding := resp.Header.Get("Content-Encoding")
	var reader io.Reader
	var err error

	switch contentEncoding {
	case "gzip":
		// 处理 gzip 压缩
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("create gzip reader failed: %v", err)
		}
		defer func() {
			if err := gzReader.Close(); err != nil {
				fmt.Printf("close gzip reader failed: %v", err)
			}
		}()
		reader = gzReader

	case "zstd":
		// 处理 zstd 压缩
		zReader, err := zstd.NewReader(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("create zstd reader failed: %v", err)
		}
		defer zReader.Close()
		reader = zReader

	case "br", "brotli":
		// 处理 brotli 压缩
		reader = brotli.NewReader(resp.Body)

	case "deflate":
		// 处理 deflate 压缩
		reader = flate.NewReader(resp.Body)
		defer func(closer io.Closer) {
			err := closer.Close()
			if err != nil {
				fmt.Printf("close deflate reader failed: %v", err)
			}
		}(reader.(io.Closer))

	case "zlib":
		// 处理 zlib 压缩
		zReader, err := zlib.NewReader(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("create zlib reader failed: %v", err)
		}
		defer func(zReader io.ReadCloser) {
			err := zReader.Close()
			if err != nil {
				fmt.Printf("close zlib reader failed: %v", err)
			}
		}(zReader)
		reader = zReader

	case "lz4":
		// 处理 LZ4 压缩
		reader = lz4.NewReader(resp.Body)

	case "snappy":
		// 处理 Snappy 压缩
		reader = snappy.NewReader(resp.Body)

	case "lzw":
		// 处理 LZW 压缩
		reader = lzw.NewReader(resp.Body, lzw.MSB, 8)
		defer func(closer io.Closer) {
			err := closer.Close()
			if err != nil {
				fmt.Printf("close lzw reader failed: %v", err)
			}
		}(reader.(io.Closer))

	case "xz":
		// 处理 XZ 压缩
		xzReader, err := xz.NewReader(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("create xz reader failed: %v", err)
		}
		reader = xzReader

	case "s2":
		// 处理 S2 压缩
		reader = s2.NewReader(resp.Body)

	default:
		// 未压缩或未知压缩方式，直接读取
		reader = resp.Body
		contentEncoding = ""
	}

	// 读取响应体
	respBody, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", fmt.Errorf("read response body failed: %v", err)
	}

	return respBody, contentEncoding, nil
}

// 将修改后的响应体返回给客户端
func writeResponse(
	c *gin.Context,
	resp *http.Response,
	modifiedBody []byte,
	contentEncoding string,
) error {
	// 如果需要，重新压缩响应体
	if contentEncoding != "" {
		var buf bytes.Buffer
		var err error

		switch contentEncoding {
		case "gzip":
			writer := gzip.NewWriter(&buf)
			if _, err = writer.Write(modifiedBody); err != nil {
				return fmt.Errorf("write gzip data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close gzip writer failed: %v", err)
			}

		case "zstd":
			writer, err := zstd.NewWriter(&buf)
			if err != nil {
				return fmt.Errorf("create zstd encoder failed: %v", err)
			}
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write zstd data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close zstd writer failed: %v", err)
			}

		case "br", "brotli":
			writer := brotli.NewWriter(&buf)
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write brotli data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close brotli writer failed: %v", err)
			}

		case "deflate":
			writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
			if err != nil {
				return fmt.Errorf("create deflate writer failed: %v", err)
			}
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write deflate data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close deflate writer failed: %v", err)
			}

		case "zlib":
			writer := zlib.NewWriter(&buf)
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write zlib data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close zlib writer failed: %v", err)
			}

		case "lz4":
			writer := lz4.NewWriter(&buf)
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write lz4 data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close lz4 writer failed: %v", err)
			}

		case "snappy":
			writer := snappy.NewBufferedWriter(&buf)
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write snappy data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close snappy writer failed: %v", err)
			}

		case "lzw":
			writer := lzw.NewWriter(&buf, lzw.MSB, 8)
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write lzw data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close lzw writer failed: %v", err)
			}

		case "xz":
			writer, err := xz.NewWriter(&buf)
			if err != nil {
				return fmt.Errorf("create xz writer failed: %v", err)
			}
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write xz data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close xz writer failed: %v", err)
			}

		case "s2":
			writer := s2.NewWriter(&buf)
			if _, err = writer.Write(modifiedBody); err != nil {
				err := writer.Close()
				if err != nil {
					return err
				}
				return fmt.Errorf("write s2 data failed: %v", err)
			}
			if err = writer.Close(); err != nil {
				return fmt.Errorf("close s2 writer failed: %v", err)
			}
		}

		modifiedBody = buf.Bytes()
	}

	// 复制上游响应的头信息到客户端响应
	for key, values := range resp.Header {
		if !strings.EqualFold(key, "Content-Length") {
			c.Writer.Header()[key] = values
		}
	}

	// 设置新的 Content-Length
	c.Writer.Header().Set("Content-Length", strconv.Itoa(len(modifiedBody)))

	// 设置响应状态码
	c.Status(resp.StatusCode)

	// 将修改后的响应体返回给客户端
	_, err := c.Writer.Write(modifiedBody)
	if err != nil {
		return fmt.Errorf("write response body failed: %v", err)
	}

	return nil
}
