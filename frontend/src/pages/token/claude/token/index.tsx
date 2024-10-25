import { useState, useEffect } from 'react';
import {
  Button,
  Card,
  Col,
  Form,
  Input,
  Modal,
  Popconfirm,
  Row,
  Space,
  Typography,
  Checkbox,
  message,
  Spin, List, Drawer
} from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import {
  DeleteOutlined,
  EditOutlined,
  ReloadOutlined
} from "@ant-design/icons";
import { siAnthropic } from 'simple-icons/icons';
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import dayjs from "dayjs";
import 'dayjs/locale/zh-cn';
import customParseFormat from 'dayjs/plugin/customParseFormat';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

import {ClaudeAccount, ClaudeToken} from '#/entity.ts';
import tokenService, { ClaudeTokenAddReq } from "@/api/services/claudeTokenService.ts";
import accountService from "@/api/services/claudeAccountService.ts";
import {
  useAddTokenMutation,
  useDeleteTokenMutation,
  useUpdateTokenMutation
} from "@/store/claudeTokenStore.ts";
import { useAddAccountMutation } from "@/store/claudeAccountStore.ts";
import CopyToClipboardInput from "@/pages/components/copy";
import formatDateTime from "@/pages/components/util";
import Chart from "@/components/chart/chart.tsx";
import useChart from "@/components/chart/useChart.ts";

dayjs.locale('zh-cn');
dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(customParseFormat);

type SearchFormFieldType = Pick<ClaudeToken, 'tokenName'>;

const LOCAL_STORAGE_KEY = 'claude_token_page_visible_columns';

export default function TokenPage() {

  const AnthropicIcon = () => {
    return (
      <div
        dangerouslySetInnerHTML={{ __html: siAnthropic.svg }}
        style={{ width: '16px', height: '16px', display: 'inline-block', verticalAlign: 'middle', fill: 'white' }}
      />
    );
  };

  const queryClient = useQueryClient();
  const [searchForm] = Form.useForm();
  const { t } = useTranslation();

  const addTokenMutation = useAddTokenMutation();
  const updateTokenMutation = useUpdateTokenMutation();
  const deleteTokenMutation = useDeleteTokenMutation();
  const addAccountMutation = useAddAccountMutation();

  const [deleteTokenId, setDeleteTokenId] = useState<number | undefined>(-1);

  const [visibleColumns, setVisibleColumns] = useState<(keyof ClaudeToken | 'operation')[]>(() => {
    const storedColumns = localStorage.getItem(LOCAL_STORAGE_KEY);
    return storedColumns
      ? JSON.parse(storedColumns)
      : ['id', 'tokenName', 'sessionToken', 'createTime', 'updateTime', 'operation'];
  });
  const [tempVisibleColumns, setTempVisibleColumns] = useState<(keyof ClaudeToken | 'operation')[]>(visibleColumns);
  const [drawerVisible, setDrawerVisible] = useState(false);

  const searchTokenName = Form.useWatch('tokenName', searchForm);

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['claudeTokens', searchTokenName],
    queryFn: () => tokenService.searchTokenList(searchTokenName),
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  });

  const [tokenModalProps, setTokenModalProps] = useState<TokenModalProps>({
    formValue: {
      tokenName: '',
      sessionToken: '',
    },
    title: 'New',
    show: false,
    onOk: (values: ClaudeTokenAddReq, callback) => {
      if (values.id) {
        updateTokenMutation.mutate(values, {
          onSuccess: () => {
            setTokenModalProps((prev) => ({ ...prev, show: false }))
            queryClient.invalidateQueries({ queryKey: ['claudeTokens'] });
          },
          onSettled: () => callback(false)
        });
      } else {
        addTokenMutation.mutate(values, {
          onSuccess: () => {
            setTokenModalProps((prev) => ({ ...prev, show: false }))
            queryClient.invalidateQueries({ queryKey: ['claudeTokens'] });
          },
          onSettled: () => callback(false)
        });
      }
    },
    onCancel: () => {
      setTokenModalProps((prev) => ({ ...prev, show: false }));
    },
  });

  const [accountModalProps, setAccountModalProps] = useState<AccountModalProps>({
    formValue: {
      userId: -1,
      tokenId: -1,
      account: '',
      status: 1,
    },
    title: 'New',
    show: false,
    isEdit: false,
    onOk: (values: ClaudeAccount, callback) => {
      callback(true);
      addAccountMutation.mutate(values, {
        onSuccess: () => {
          setAccountModalProps((prev) => ({ ...prev, show: false }))
          queryClient.invalidateQueries({ queryKey: ['claudeTokens'] });
        },
        onSettled: () => callback(false)
      });
    },
    onCancel: () => {
      setAccountModalProps((prev) => ({ ...prev, show: false }));
    },
  });

  const [accountInfoModalProps, setAccountInfoModalProps] = useState<AccountInfoModalProps>({
    tokenId: -1,
    show: false,
    onOk: () => {
      setAccountInfoModalProps((prev) => ({ ...prev, show: false }));
    },
  });

  const [chatTokenId, setChatTokenId] = useState<number | undefined>(-1);

  function handleQuickLogin(record: ClaudeToken) {
    let id = record.id ? record.id : -5;
    accountService.chatAuthAccount(5, id)
      .then((res) => {
        const {loginUrl} = res;
        if (loginUrl) {
          window.open(loginUrl)
        } else {
          message.error('Failed to get login url')
        }
      })
      .catch((err) => {
        console.log(err)
        message.error('Failed to get login url')
      })
      .finally(() => {
        setChatTokenId(undefined)
      })
  }

  const columns: ColumnsType<ClaudeToken> = [
    {
      title: t('token.id'),
      key: 'id',
      dataIndex: 'id',
      ellipsis: true,
      align: 'center',
      render: (text) => (
        <Typography.Text style={{ maxWidth: 200 }} ellipsis={true}>
          {text}
        </Typography.Text>
      )
    },
    {
      title: t('token.tokenName'),
      key: 'tokenName',
      dataIndex: 'tokenName',
      align: 'center',
      ellipsis: true,
      width: 200,
      render: (text) => (
        <CopyToClipboardInput text={text} showTooltip={true} />
      )
    },
    {
      title: t('token.sessionToken'),
      key: 'sessionToken',
      dataIndex: 'sessionToken',
      align: 'center',
      width: 250,
      render: (text) => (
        <CopyToClipboardInput text={text} />
      ),
    },
    {
      title: t("token.createTime"),
      key: 'createTime',
      dataIndex: 'createTime',
      align: 'center',
      width: 200,
      render: (text) => formatDateTime(text),
    },
    {
      title: t("token.updateTime"),
      key: 'updateTime',
      dataIndex: 'updateTime',
      align: 'center',
      width: 200,
      render: (text) => formatDateTime(text),
    },
    {
      title: t('token.action'),
      key: 'operation',
      align: 'center',
      render: (_, record) => (
        <Button.Group>
          <Button
            icon={<AnthropicIcon />}
            type={"primary"}
            onClick={() => handleQuickLogin(record)}
            loading={chatTokenId === record.id}
            style={{ backgroundColor: '#007bff', borderColor: '#007bff', color: 'white' }}
          >Chat</Button>
          <Button onClick={() => onEdit(record)} icon={<EditOutlined />} type="primary" />
          <Popconfirm title={t('common.deleteConfirm')} okText={t('common.yes')} cancelText={t('common.no')} placement="left" onConfirm={() => {
            setDeleteTokenId(record.id);
            deleteTokenMutation.mutate(record.id, {
              onSuccess: () => {
                setDeleteTokenId(undefined)
                queryClient.invalidateQueries({ queryKey: ['claudeTokens'] });
              }
            })
          }}>
            <Button icon={<DeleteOutlined />} type="primary" loading={deleteTokenId === record.id} danger />
          </Popconfirm>
        </Button.Group>
      ),
    },
  ];

  useEffect(() => {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(visibleColumns));
  }, [visibleColumns]);

  const showDrawer = () => {
    setDrawerVisible(true);
  };

  const onDrawerClose = () => {
    setDrawerVisible(false);
    setTempVisibleColumns(visibleColumns);
  };

  const applyColumnVisibility = () => {
    setVisibleColumns(tempVisibleColumns);
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(tempVisibleColumns));
    setDrawerVisible(false);
  };

  const selectAll = () => {
    const allColumnKeys = columns.map(col => col.key as keyof ClaudeToken | 'operation');
    setTempVisibleColumns(allColumnKeys);
  };

  const deselectAll = () => {
    setTempVisibleColumns([]);
  };

  const visibleColumnsConfig = columns.filter(col =>
    col.key && visibleColumns.includes(col.key as keyof ClaudeToken | 'operation')
  );

  const onSearchFormReset = () => {
    searchForm.resetFields();
  };

  const onCreate = () => {
    setTokenModalProps((prev) => ({
      ...prev,
      show: true,
      title: t('token.createNew'),
      formValue: {
        id: undefined,
        tokenName: '',
        sessionToken: '',
      },
    }));
  };

  const onEdit = (record: ClaudeToken) => {
    setTokenModalProps((prev) => ({
      ...prev,
      show: true,
      title: t('token.edit'),
      formValue: {
        id: record.id,
        tokenName: record.tokenName,
        sessionToken: record.sessionToken,
      },
    }));
  };

  const handleRefresh = () => {
    refetch();
    message.success(t('common.dataRefreshed'));
  };

  return (
    <Space direction="vertical" size="large" className="w-full">
      <Card>
        <Form form={searchForm}>
          <Row gutter={[16, 16]}>
            <Col span={6} lg={6}>
              <Form.Item<SearchFormFieldType> label={t('token.tokenName')} name="tokenName" className="!mb-0">
                <Input />
              </Form.Item>
            </Col>
            <Col span={18} lg={18}>
              <div className="flex justify-end">
                <Space>
                  <Button onClick={onSearchFormReset}>{t('token.reset')}</Button>
                  <Button icon={<ReloadOutlined />} onClick={handleRefresh}>
                    {t('common.refresh')}
                  </Button>
                </Space>
              </div>
            </Col>
          </Row>
        </Form>
      </Card>

      <Card
        title={t("token.accountList")}
        extra={
          <Space>
            <Button onClick={showDrawer}>
              {t("token.adjustDisplay")}
            </Button>
            <Button type="primary" onClick={onCreate}>
              {t("token.createNew")}
            </Button>
          </Space>
        }
      >
        <Table
          rowKey="id"
          size="small"
          scroll={{ x: 'max-content' }}
          pagination={{ pageSize: 10 }}
          columns={visibleColumnsConfig}
          dataSource={data}
          loading={isLoading}
        />
      </Card>

      <Drawer
        title={t("token.selectColumns")}
        placement="right"
        onClose={onDrawerClose}
        open={drawerVisible}
        width={260} // 可以稍微减小宽度，因为我们去掉了额外的描述文本
        extra={
          <Space>
            <Button onClick={applyColumnVisibility} type="primary">
              {t('common.apply')}
            </Button>
          </Space>
        }
      >
        <div style={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
          <div style={{ marginBottom: '16px' }}>
            <Space>
              <Button
                size="small" // 增大按钮尺寸
                type="default" // 使用默认类型，避免过于鲜艳
                onClick={selectAll}
                style={{
                  width: '100px', // 设置按钮宽度
                  height: '40px',  // 设置按钮高度
                  borderRadius: '8px', // 圆角调整
                  backgroundColor: '#e6f7ff', // 柔和的蓝色背景
                  borderColor: '#91d5ff', // 边框颜色
                  color: '#1890ff', // 文字颜色
                  transition: 'background-color 0.3s, border-color 0.3s, color 0.3s',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = '#bae7ff';
                  e.currentTarget.style.borderColor = '#40a9ff';
                  e.currentTarget.style.color = '#096dd9';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = '#e6f7ff';
                  e.currentTarget.style.borderColor = '#91d5ff';
                  e.currentTarget.style.color = '#1890ff';
                }}
              >
                {t('common.selectAll')}
              </Button>

              <Button
                size="small" // 增大按钮尺寸
                type="default" // 使用默认类型，避免过于鲜艳
                onClick={deselectAll}
                style={{
                  width: '100px', // 设置按钮宽度
                  height: '40px',  // 设置按钮高度
                  borderRadius: '8px', // 圆角调整
                  backgroundColor: '#fff1f0', // 柔和的红色背景
                  borderColor: '#ffa39e', // 边框颜色
                  color: '#ff4d4f', // 文字颜色
                  transition: 'background-color 0.3s, border-color 0.3s, color 0.3s',
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.backgroundColor = '#ffa39e';
                  e.currentTarget.style.borderColor = '#ff7875';
                  e.currentTarget.style.color = '#a8071a';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.backgroundColor = '#fff1f0';
                  e.currentTarget.style.borderColor = '#ffa39e';
                  e.currentTarget.style.color = '#ff4d4f';
                }}
              >
                {t('common.deselectAll')}
              </Button>
            </Space>
          </div>
          <List
            style={{
              flexGrow: 1,
              overflowY: 'auto',
            }}
            dataSource={columns}
            renderItem={col => (
              <List.Item style={{ border: 'none', padding: '8px 0' }}> {/* 移除边框 */}
                <Checkbox
                  checked={tempVisibleColumns.includes(col.key as keyof ClaudeToken | 'operation')}
                  onChange={(e) => {
                    const checked = e.target.checked;
                    if (checked) {
                      setTempVisibleColumns([...tempVisibleColumns, col.key as keyof ClaudeToken | 'operation']);
                    } else {
                      setTempVisibleColumns(tempVisibleColumns.filter(k => k !== col.key));
                    }
                  }}
                  style={{ width: '100%' }} // 让 Checkbox 占满整行
                >
                  {typeof col.title === 'function' ? col.title({}) : col.title}
                </Checkbox>
              </List.Item>
            )}
          />
        </div>
      </Drawer>

      <TokenModal {...tokenModalProps} />
      <AccountModal {...accountModalProps} />
      <AccountInfoModal {...accountInfoModalProps} />
    </Space>
  );
}

export type AccountModalProps = {
  formValue: ClaudeAccount;
  title: string;
  show: boolean;
  isEdit: boolean;
  onOk: (values: ClaudeAccount, callback: any) => void;
  onCancel: VoidFunction;
}

export const AccountModal = ({ title, show, isEdit, formValue, onOk, onCancel }: AccountModalProps) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (show) {
      form.setFieldsValue(formValue);
    } else {
      form.resetFields();
    }
  }, [formValue, show, form]);

  const onModalOk = () => {
    form.validateFields().then((values) => {
      const formattedValues = {
        ...values,
      };
      setLoading(true);
      onOk(formattedValues, () => setLoading(false));
    }).catch(error => {
      console.error('Validation error:', error);
    });
  };

  return (
    <Modal
      title={title}
      open={show}
      onOk={onModalOk}
      onCancel={onCancel}
      okButtonProps={{
        loading: loading,
      }}
      destroyOnClose={true}
    >
      <Form
        form={form}
        layout="vertical"
      >
        <Form.Item<ClaudeAccount> name="id" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<ClaudeAccount> name="userId" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<ClaudeAccount> name="tokenId" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<ClaudeAccount> label="ClaudeAccount" name="account" required>
          <Input readOnly={isEdit} disabled={isEdit} autoComplete="off"/>
        </Form.Item>
      </Form>
    </Modal>
  );
}

type TokenModalProps = {
  formValue: ClaudeTokenAddReq;
  title: string;
  show: boolean;
  onOk: (values: ClaudeTokenAddReq, setLoading: any) => void;
  onCancel: VoidFunction;
};

function TokenModal({title, show, formValue, onOk, onCancel}: TokenModalProps) {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation();

  useEffect(() => {
    if (show) {
      form.setFieldsValue(formValue);
    } else {
      form.resetFields();
    }
  }, [show, formValue, form]);

  const onModalOk = () => {
    form.validateFields().then((values) => {
      setLoading(true);
      onOk(values, setLoading);
    });
  };

  return (
    <Modal
      title={title}
      open={show}
      onOk={onModalOk}
      onCancel={onCancel}
      okButtonProps={{ loading: loading }}
      destroyOnClose={true}
    >
      <Form
        form={form}
        layout="vertical"
      >
        <Form.Item<ClaudeTokenAddReq> name="id" hidden>
          <Input />
        </Form.Item>
        <Form.Item<ClaudeTokenAddReq> label={t("token.tokenName")} name="tokenName" required>
          <Input autoComplete="off" />
        </Form.Item>
        <Form.Item<ClaudeTokenAddReq> label={t("token.sessionToken")} name="sessionToken" required>
          <Input autoComplete="off" />
        </Form.Item>
      </Form>
    </Modal>
  );
}

type AccountInfoModalProps = {
  tokenId: number
  onOk: VoidFunction
  show: boolean;
}

const AccountInfoModal = ({tokenId, onOk, show}: AccountInfoModalProps) => {
  const {data: statistic, isLoading} = useQuery({
    queryKey: ['claudeTokenStatistic', tokenId],
    queryFn: () => accountService.getAccountStatistic(tokenId),
    enabled: show,
  })

  const {t} = useTranslation()

  let chartOptions = useChart({
    legend: {
      horizontalAlign: 'center',
    },
    stroke: {
      show: true,
    },
    dataLabels: {
      enabled: true,
      dropShadow: {
        enabled: false,
      },
    },
    xaxis: {
      categories: statistic?.categories || [],
    },
    tooltip: {
      fillSeriesColor: false,
    },
    plotOptions: {
      pie: {
        donut: {
          labels: {
            show: false,
          },
        },
      },
    },
  });

  return (
    <Modal title={t('token.statistic')} open={show} onOk={onOk} closable={false} onCancel={onOk}>
      <Spin spinning={isLoading} tip={t("token.queryingInfo")}>
        <Chart type="bar" series={statistic?.series || []} options={chartOptions} height={320}/>
      </Spin>
    </Modal>
  )
}
