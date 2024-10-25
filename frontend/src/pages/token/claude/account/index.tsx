import {
  Button,
  Card,
  Checkbox,
  Col,
  Drawer,
  Form,
  Input,
  List,
  message,
  Row,
  Space,
  Tooltip,
} from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import {useEffect, useState} from 'react';

import {ClaudeAccount} from '#/entity.ts';
import {
  CheckCircleOutlined, CloseCircleOutlined, EditOutlined,
  ReloadOutlined,
} from "@ant-design/icons";
import { siAnthropic } from 'simple-icons/icons';
import { useQuery, useQueryClient } from "@tanstack/react-query";
import {useSearchParams} from "@/router/hooks";
import accountService from "@/api/services/claudeAccountService.ts";
import {
  useUpdateAccountMutation
} from "@/store/claudeAccountStore.ts";
import {AccountModal, AccountModalProps} from "@/pages/token/claude/token";
import {useTranslation} from "react-i18next";
import formatDateTime from "@/pages/components/util";
import CopyToClipboardInput from "@/pages/components/copy";

type SearchFormFieldType = {
  tokenId?: number;
};

const LOCAL_STORAGE_KEY = 'claude_share_page_visible_columns';

export default function SharePage() {

  const AnthropicIcon = () => {
    return (
      <div
        dangerouslySetInnerHTML={{ __html: siAnthropic.svg }}
        style={{ width: '16px', height: '16px', display: 'inline-block', verticalAlign: 'middle', fill: 'white' }}
      />
    );
  };

  const queryClient = useQueryClient();
  const updateShareMutation = useUpdateAccountMutation()

  const [visibleColumns, setVisibleColumns] = useState<(keyof ClaudeAccount | 'operation')[]>(() => {
    const storedColumns = localStorage.getItem(LOCAL_STORAGE_KEY);
    return storedColumns
      ? JSON.parse(storedColumns)
      : ['tokenId', 'account', 'status', 'createTime', 'updateTime', 'operation'];
  });
  const [tempVisibleColumns, setTempVisibleColumns] = useState<(keyof ClaudeAccount | 'operation')[]>(visibleColumns);
  const [drawerVisible, setDrawerVisible] = useState(false);

  const {t} = useTranslation()

  const params = useSearchParams();
  const [searchForm] = Form.useForm();
  const tokenId = Form.useWatch('tokenId', searchForm);
  const [shareModalProps, setShareModalProps] = useState<AccountModalProps>({
    formValue: {
      userId: -1,
      tokenId: -1,
      account: '',
      status: 1,
    },
    title: t('token.edit'),
    show: false,
    isEdit: false,
    onOk: (values: ClaudeAccount) => {
      console.log(values)
      setShareModalProps((prev) => ({...prev, show: false}));
    },
    onCancel: () => {
      setShareModalProps((prev) => ({...prev, show: false}));
    },
  });
  const [chatAccountId, setChatAccountId] = useState<number | undefined>(-1);

  useEffect(() => {
    searchForm.setFieldValue('tokenId', params.get('tokenId'))
  }, [params]);

  // Add this effect to invalidate query cache when the component mounts
  useEffect(() => {
    queryClient.invalidateQueries({ queryKey: ['claudeAccounts'] });
  }, [queryClient]);

  function handleQuickLogin(record: ClaudeAccount) {
    let id = record.id ? record.id : -4;
    accountService.chatAuthAccount(4, id)
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
        setChatAccountId(undefined)
      })
  }

  const columns: ColumnsType<ClaudeAccount> = [
    {
      title: t('token.tokenId'),
      key: 'tokenId',
      dataIndex: 'tokenId',
      align: 'center',
      width: 80
    },
    {
      title: t('token.user.claude'),
      key: 'account',
      dataIndex: 'account',
      align: 'center',
      width: 120,
      render: (text) => (
        <CopyToClipboardInput text={text} showTooltip={true} />
      )
    },
    {
      title: t('token.accountStatus'),
      key: 'status',
      dataIndex: 'status',
      align: 'center',
      render: (status) => {
        if (status === 0) {
          return <Tooltip title={t('token.disable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (status === 1) {
          return <Tooltip title={t('token.normal')}><CheckCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        }
      },
    },
    {
      title: t('token.createTime'),
      key: 'createTime',
      dataIndex: 'createTime',
      align: 'center',
      width: 200,
      render: (text) => formatDateTime(text),
    },
    {
      title: t('token.updateTime'),
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
      render: (_,record) => (
        <Button.Group>
          <Button
            icon={<AnthropicIcon />}
            type={"primary"}
            onClick={() => handleQuickLogin(record)}
            loading={chatAccountId === record.id}
            style={{ backgroundColor: '#007bff', borderColor: '#007bff', color: 'white' }}
            disabled={record.status !== 1}
          >
            Chat
          </Button>
          <Button icon={<EditOutlined />} type={"primary"} onClick={() => onEdit(record)}></Button>
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
    const allColumnKeys = columns.map(col => col.key as keyof ClaudeAccount | 'operation');
    setTempVisibleColumns(allColumnKeys);
  };

  const deselectAll = () => {
    setTempVisibleColumns([]);
  };

  const visibleColumnsConfig = columns.filter(col =>
    col.key && visibleColumns.includes(col.key as keyof ClaudeAccount | 'operation')
  );

  const onEdit = (record: ClaudeAccount) => {
    console.log(record)
    setShareModalProps({
      formValue: record,
      title: t('token.edit'),
      show: true,
      isEdit: true,
      onOk: (values: ClaudeAccount, callback) => {
        updateShareMutation.mutate(values, {
          onSuccess: () => {
            setShareModalProps((prev) => ({...prev, show: false}));
            queryClient.invalidateQueries({ queryKey: ['claudeAccounts'] });
          },
          onSettled: () => callback(false)
        })
      },
      onCancel: () => {
        setShareModalProps((prev) => ({...prev, show: false}));
      },
    })
  }

  const { data, refetch, isLoading } = useQuery({
    queryKey: ['claudeAccounts', tokenId],
    queryFn: () => {
      let tokenIdNum = parseInt(tokenId as any);
      return accountService.searchAccount(tokenIdNum);
    }
  })

  // Add this effect to refetch data when the component mounts or when tokenId changes
  useEffect(() => {
    refetch();
  }, [refetch, tokenId]);

  const onSearchFormReset = () => {
    searchForm.resetFields();
    refetch();
  };

  const handleRefresh = () => {
    refetch();
    message.success(t('common.dataRefreshed'));
  };

  return (
    <Space direction="vertical" size="large" className="w-full">
      <Card>
        <Form form={searchForm} >
          <Row gutter={[16, 16]}>
            <Col span={3} lg={3}>
              <Form.Item<SearchFormFieldType> label={t('token.tokenId')} name="tokenId" className="!mb-0">
                <Input />
              </Form.Item>
            </Col>
            <Col span={21} lg={21}>
              <div className="flex justify-end">
                <Button onClick={onSearchFormReset}>{t('token.reset')}</Button>
                <Button
                  icon={<ReloadOutlined />}
                  onClick={handleRefresh}
                  loading={isLoading}
                >
                  {t('common.refresh')}
                </Button>
              </div>
            </Col>
          </Row>
        </Form>
      </Card>

      <Card
        title={t('token.shareList')}
        extra={
          <Space>
            <Button onClick={showDrawer}>
              {t("token.adjustDisplay")}
            </Button>
          </Space>
        }
      >
        <Table
          rowKey={record => record.id + record.account}
          size="small"
          scroll={{ x: 'max-content' }}
          pagination={{ pageSize: 10 }}
          columns={visibleColumnsConfig}
          dataSource={data}
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
                  checked={tempVisibleColumns.includes(col.key as keyof ClaudeAccount | 'operation')}
                  onChange={(e) => {
                    const checked = e.target.checked;
                    if (checked) {
                      setTempVisibleColumns([...tempVisibleColumns, col.key as keyof ClaudeAccount | 'operation']);
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

      <AccountModal {...shareModalProps}/>
    </Space>
  );
}
