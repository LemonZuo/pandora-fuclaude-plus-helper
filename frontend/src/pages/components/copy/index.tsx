import React from 'react';
import {Input, Button, message, Tooltip} from 'antd';
import {CopyOutlined} from '@ant-design/icons';

// 定义组件的 props 类型
interface CopyToClipboardInputProps {
  text: string;  // 指定 text 为 string 类型
  showTooltip?: boolean;
}

const CopyToClipboardInput: React.FC<CopyToClipboardInputProps> = ({text, showTooltip = false}) => {
  const handleCopy = async (text: string) => {
    try {
      await navigator.clipboard.writeText(text);
      message.success('复制成功');
    } catch (err) {
      message.error('复制失败');
    }
  };

  return (
    <div style={{display: 'flex', alignItems: 'center', gap: '5px'}}>
      {showTooltip ? (
        <Tooltip title={text} placement="top">
          <Input value={text} readOnly style={{flex: 1}}/>
        </Tooltip>
      ) : (
        <Input value={text} readOnly style={{flex: 1}}/>
      )}
      <Button icon={<CopyOutlined/>} onClick={() => handleCopy(text)}/>
    </div>
  );
};

export default CopyToClipboardInput;
