import {
  Button,
  Card, Checkbox, CheckboxOptionType,
  Col,
  Form,
  Input,
  message,
  Popover,
  Row,
  Space, Tooltip,
} from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import {useEffect, useState} from 'react';

import {OpenaiAccount} from '#/entity.ts';
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  EditOutlined,
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
      : ['tokenId', 'account', 'status', 'expirationTime',
        'shareToken', 'gpt35Limit', 'gpt4Limit', 'showConversations',
        'temporaryChat', 'expireAt', 'createTime', 'updateTime', 'operation'];
  });
  const [tempVisibleColumns, setTempVisibleColumns] = useState<(keyof OpenaiAccount | 'operation')[]>(visibleColumns);
  const [popoverVisible, setPopoverVisible] = useState(false);

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
      width: 120
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
    {
      title: t('token.gpt35Limit'),
      key: 'gpt35Limit',
      dataIndex: 'gpt35Limit',
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

  const handleVisibilityChange = (checkedValues: (keyof OpenaiAccount | 'operation')[]) => {
    setTempVisibleColumns(checkedValues);
  };

  const applyColumnVisibility = () => {
    setVisibleColumns(tempVisibleColumns);
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(tempVisibleColumns));
    setPopoverVisible(false);
  };

  const columnVisibilityContent = (
    <div style={{ maxWidth: 120 }}>
      <Checkbox.Group
        options={columns.map(col => ({ label: col.title, value: col.key })) as CheckboxOptionType<keyof OpenaiAccount | "operation">[]}
        value={tempVisibleColumns}
        onChange={handleVisibilityChange}
        style={{display: 'block'}}
      />
      <div style={{ marginTop: 8, textAlign: 'right' }}>
        <Button size="small" type="primary" onClick={applyColumnVisibility}>
          {t('common.apply')}
        </Button>
      </div>
    </div>
  );

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
            <Popover
              content={columnVisibilityContent}
              title={t("token.selectColumns")}
              trigger="click"
              open={popoverVisible}
              onOpenChange={setPopoverVisible}
            >
              <Button>
                {t("token.adjustDisplay")}
              </Button>
            </Popover>
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
      <AccountModal {...shareModalProps}/>
    </Space>
  );
}
