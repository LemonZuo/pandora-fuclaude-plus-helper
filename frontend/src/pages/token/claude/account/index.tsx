import {
  Button,
  Card,
  Checkbox,
  CheckboxOptionType,
  Col,
  Form,
  Input,
  message,
  Popover,
  Row,
  Space,
  Tooltip,
} from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import {useEffect, useState} from 'react';

import {ClaudeAccount} from '#/entity.ts';
import {
  CheckCircleOutlined, CloseCircleOutlined,
  EditOutlined,
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

  const handleVisibilityChange = (checkedValues: (keyof ClaudeAccount | 'operation')[]) => {
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
        options={columns.map(col => ({ label: col.title, value: col.key })) as CheckboxOptionType<keyof ClaudeAccount | "operation">[]}
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
        />
      </Card>
      <AccountModal {...shareModalProps}/>
    </Space>
  );
}
