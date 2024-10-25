import {
  Button,
  Card, Checkbox, Col, Drawer,
  Form,
  Input,
  List,
  message,
  Row,
  Space, Tooltip,
} from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import {useEffect, useState} from 'react';

import {OpenaiAccount} from '#/entity.ts';
import {
  CheckCircleOutlined, CloseCircleOutlined, EditOutlined,
  ExclamationCircleOutlined,
  MinusCircleOutlined,
  OpenAIFilled,
  ReloadOutlined,
} from "@ant-design/icons";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import {useSearchParams} from "@/router/hooks";
import accountService from "@/api/services/accountService.ts";
import {
  useUpdateAccountMutation
} from "@/store/accountStore.ts";
import {AccountModal, AccountModalProps} from "@/pages/token/openai/token";
import {useTranslation} from "react-i18next";
import CopyToClipboardInput from "@/pages/components/copy";
import formatDateTime from "@/pages/components/util";

type SearchFormFieldType = {
  tokenId?: number;
};

const LOCAL_STORAGE_KEY = 'openai_share_page_visible_columns';

export default function SharePage() {
  const queryClient = useQueryClient();
  const updateShareMutation = useUpdateAccountMutation()

  const [visibleColumns, setVisibleColumns] = useState<(keyof OpenaiAccount | 'operation')[]>(() => {
    const storedColumns = localStorage.getItem(LOCAL_STORAGE_KEY);
    return storedColumns
      ? JSON.parse(storedColumns)
      : ['tokenId', 'account', 'status', 'expirationTime','shareToken',
        'gpt35Limit', 'gpt4Limit', 'gpt4oLimit', 'gpt4oMiniLimit', 'o1Limit', 'o1MiniLimit',
        'showConversations','temporaryChat', 'expireAt', 'createTime', 'updateTime', 'operation'];
  });
  const [tempVisibleColumns, setTempVisibleColumns] = useState<(keyof OpenaiAccount | 'operation')[]>(visibleColumns);
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
      expirationTime: '',
      gpt35Limit: -1,
      gpt4Limit: -1,
      gpt4oLimit: -1,
      gpt4oMiniLimit: -1,
      o1Limit: -1,
      o1MiniLimit: -1,
      showConversations: 0,
      temporaryChat: 0,
    },
    title: t('token.edit'),
    show: false,
    isEdit: false,
    onOk: (values: OpenaiAccount) => {
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

  function handleQuickLogin(record: OpenaiAccount) {
    let id = record.id ? record.id : -2;
    accountService.chatAuthAccount(2, id)
      .then((res) => {
        const {loginUrl} = res;
        if (loginUrl) {
          window.open(loginUrl)
        } else {
          message.error('Failed to get login url').then(r => console.log(r))
        }
      })
      .catch((err) => {
        console.log(err)
        message.error('Failed to get login url').then(r => console.log(r))
      })
      .finally(() => {
        setChatAccountId(undefined)
      })
  }

  const columns: ColumnsType<OpenaiAccount> = [
    {
      title: t('token.tokenId'),
      key: 'tokenId',
      dataIndex: 'tokenId',
      align: 'center',
      width: 80
    },
    {
      title: t('token.user.openai'),
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
      title: 'ShareToken',
      key: 'shareToken',
      dataIndex: 'shareToken',
      align: 'center',
      render: (text) => (
        <CopyToClipboardInput text={text}/>
      ),
    },
    // {
    //   title: t('token.gpt35Limit'),
    //   key: 'gpt35Limit',
    //   dataIndex: 'gpt35Limit',
    //   align: 'center',
    //   width: 120,
    //   render: (count) => {
    //     if (count === 0) {
    //       return <Tooltip title={t('token.notAvailable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
    //     } else if (count < 0) {
    //       return <Tooltip title={t('token.unlimitedTimes')}><MinusCircleOutlined style={{ color: 'green' }} /></Tooltip>;
    //     } else {
    //       return (
    //         <Tooltip title={`${t('token.limitedTimes')}:${count}`}>
    //           <ExclamationCircleOutlined style={{ color: 'orange' }} />
    //         </Tooltip>
    //       );
    //     }
    //   },
    // },
    {
      title: t('token.gpt4Limit'),
      key: 'gpt4Limit',
      dataIndex: 'gpt4Limit',
      align: 'center',
      width: 120,
      render: (count) => {
        if (count === 0) {
          return <Tooltip title={t('token.notAvailable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (count < 0) {
          return <Tooltip title={t('token.unlimitedTimes')}><MinusCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        } else {
          return (
            <Tooltip title={`${t('token.limitedTimes')}:${count}`}>
              <ExclamationCircleOutlined style={{ color: 'orange' }} />
            </Tooltip>
          );
        }
      },
    },
    {
      title: t('token.gpt4oLimit'),
      key: 'gpt4oLimit',
      dataIndex: 'gpt4oLimit',
      align: 'center',
      width: 120,
      render: (count) => {
        if (count === 0) {
          return <Tooltip title={t('token.notAvailable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (count < 0) {
          return <Tooltip title={t('token.unlimitedTimes')}><MinusCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        } else {
          return (
            <Tooltip title={`${t('token.limitedTimes')}:${count}`}>
              <ExclamationCircleOutlined style={{ color: 'orange' }} />
            </Tooltip>
          );
        }
      },
    },
    {
      title: t('token.gpt4oMiniLimit'),
      key: 'gpt4oMiniLimit',
      dataIndex: 'gpt4oMiniLimit',
      align: 'center',
      width: 130,
      render: (count) => {
        if (count === 0) {
          return <Tooltip title={t('token.notAvailable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (count < 0) {
          return <Tooltip title={t('token.unlimitedTimes')}><MinusCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        } else {
          return (
            <Tooltip title={`${t('token.limitedTimes')}:${count}`}>
              <ExclamationCircleOutlined style={{ color: 'orange' }} />
            </Tooltip>
          );
        }
      },
    },
    {
      title: t('token.o1Limit'),
      key: 'o1Limit',
      dataIndex: 'o1Limit',
      align: 'center',
      width: 120,
      render: (count) => {
        if (count === 0) {
          return <Tooltip title={t('token.notAvailable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (count < 0) {
          return <Tooltip title={t('token.unlimitedTimes')}><MinusCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        } else {
          return (
            <Tooltip title={`${t('token.limitedTimes')}:${count}`}>
              <ExclamationCircleOutlined style={{ color: 'orange' }} />
            </Tooltip>
          );
        }
      },
    },
    {
      title: t('token.o1MiniLimit'),
      key: 'o1MiniLimit',
      dataIndex: 'o1MiniLimit',
      align: 'center',
      width: 120,
      render: (count) => {
        if (count === 0) {
          return <Tooltip title={t('token.notAvailable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (count < 0) {
          return <Tooltip title={t('token.unlimitedTimes')}><MinusCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        } else {
          return (
            <Tooltip title={`${t('token.limitedTimes')}:${count}`}>
              <ExclamationCircleOutlined style={{ color: 'orange' }} />
            </Tooltip>
          );
        }
      },
    },
    {
      title: t('token.showConversations'),
      key: 'showConversations',
      dataIndex: 'showConversations',
      align: 'center',
      width: 120,
      render: (text) => {
        if (text === 1) {
          return (
            <Tooltip title={t('common.yes')}>
              <CheckCircleOutlined style={{ color: 'orange' }} />
            </Tooltip>
          );
        } else {
          return (
            <Tooltip title={t('common.no')}>
              <CloseCircleOutlined style={{ color: 'green' }} />
            </Tooltip>
          );
        }
      },
    },
    {
      title: t('token.temporaryChat'),
      key: 'temporaryChat',
      dataIndex: 'temporaryChat',
      align: 'center',
      width: 120,
      render: (text) => {
        if (text === 1) {
          return (
            <Tooltip title={t('common.yes')}>
              <CheckCircleOutlined style={{ color: 'red' }} />
            </Tooltip>
          );
        } else {
          return (
            <Tooltip title={t('common.no')}>
              <CloseCircleOutlined style={{ color: 'green' }} />
            </Tooltip>
          );
        }
      },
    },
    {
      title: t('token.expireAt'),
      key: 'expireAt',
      dataIndex: 'expireAt',
      align: 'center',
      width: 200,
      render: (text) => formatDateTime(text),
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
            icon={<OpenAIFilled />}
            type={"primary"}
            onClick={() => handleQuickLogin(record)}
            loading={chatAccountId === record.id}
            style={{ backgroundColor: '#007bff', borderColor: '#007bff', color: 'white' }}
            disabled={record.status !== 1}
          >Chat</Button>
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
    const allColumnKeys = columns.map(col => col.key as keyof OpenaiAccount | 'operation');
    setTempVisibleColumns(allColumnKeys);
  };

  const deselectAll = () => {
    setTempVisibleColumns([]);
  };

  const visibleColumnsConfig = columns.filter(col =>
    col.key && visibleColumns.includes(col.key as keyof OpenaiAccount | 'operation')
  );

  const onEdit = (record: OpenaiAccount) => {
    setShareModalProps({
      formValue: record,
      title: t('token.edit'),
      show: true,
      isEdit: true,
      onOk: (values: OpenaiAccount, callback) => {
        values.gpt35Limit = parseInt(values.gpt35Limit as any);
        values.gpt4Limit = parseInt(values.gpt4Limit as any);
        values.gpt4oLimit = parseInt(values.gpt4oLimit as any);
        values.gpt4oMiniLimit = parseInt(values.gpt4oMiniLimit as any);
        values.o1Limit = parseInt(values.o1Limit as any);
        values.o1MiniLimit = parseInt(values.o1MiniLimit as any);
        updateShareMutation.mutate(values, {
          onSuccess: () => {
            setShareModalProps((prev) => ({...prev, show: false}));
            queryClient.invalidateQueries({ queryKey: ['openaiAccounts'] });
          },
          onSettled: () => callback(false)
        })
      },
      onCancel: () => {
        setShareModalProps((prev) => ({...prev, show: false}));
      },
    })
  }

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['openaiAccounts', tokenId],
    queryFn: () => {
      let tokenIdNum = parseInt(tokenId as any);
      return accountService.searchAccount(tokenIdNum);
    },
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  })

  const onSearchFormReset = () => {
    searchForm.resetFields();
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
                  checked={tempVisibleColumns.includes(col.key as keyof OpenaiAccount | 'operation')}
                  onChange={(e) => {
                    const checked = e.target.checked;
                    if (checked) {
                      setTempVisibleColumns([...tempVisibleColumns, col.key as keyof OpenaiAccount | 'operation']);
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
