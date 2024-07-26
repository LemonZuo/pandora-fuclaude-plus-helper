import React, { useState, useEffect } from 'react';
import {
  Button,
  Card,
  Col,
  DatePicker,
  Form,
  Input,
  Modal,
  Popconfirm,
  Row,
  Select,
  Space,
  Typography,
  Checkbox, Popover, Tooltip, message
} from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import {
  CheckCircleOutlined,
  CloseCircleOutlined,
  DeleteOutlined,
  EditOutlined,
  ReloadOutlined
} from "@ant-design/icons";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import dayjs from "dayjs";
import 'dayjs/locale/zh-cn';
import customParseFormat from 'dayjs/plugin/customParseFormat';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

import {
  useAddUserMutation,
  useDeleteUserMutation,
  useUpdateUserMutation
} from "@/store/userStore.ts";
import CopyToClipboardInput from "@/pages/components/copy";
import formatDateTime from "@/pages/components/util";
import { User } from "#/entity.ts";
import userService, { UserAddReq } from "@/api/services/userService.ts";
import tokenService from "@/api/services/tokenService.ts";
import claudeTokenService from "@/api/services/claudeTokenService.ts";

dayjs.locale('zh-cn');
dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(customParseFormat);

type SearchFormFieldType = Pick<User, 'uniqueName'>;

const { Option } = Select;

const LOCAL_STORAGE_KEY = 'user_page_visible_columns';

export default function UserPage() {
  const queryClient = useQueryClient();
  const [searchForm] = Form.useForm();
  const { t } = useTranslation();

  const addUserMutation = useAddUserMutation();
  const updateUserMutation = useUpdateUserMutation();
  const deleteUserMutation = useDeleteUserMutation();

  const [deleteUserId, setDeleteUserId] = useState<number | undefined>(-1);

  const [visibleColumns, setVisibleColumns] = useState<(keyof User | 'operation')[]>(() => {
    const storedColumns = localStorage.getItem(LOCAL_STORAGE_KEY);
    return storedColumns
      ? JSON.parse(storedColumns)
      : ['id', 'uniqueName', 'password', 'enable',
        'openai', 'openaiToken',
        'claude',
        'expireAt', 'createTime', 'updateTime', 'operation'];
  });
  const [tempVisibleColumns, setTempVisibleColumns] = useState<(keyof User | 'operation')[]>(visibleColumns);
  const [popoverVisible, setPopoverVisible] = useState(false);

  const uniqueName = Form.useWatch('uniqueName', searchForm);

  const [userModalProps, setUserModalProps] = useState<UserModalProps>({
    formValue: {
      id: -1,
      uniqueName: '',
      password: '',
      enable: 1,
      openai: 1,
      claude: 1,
      expirationTime: '',
    },
    title: 'New',
    show: false,
    isEdit: false,
    onOk: (values: UserAddReq, callback) => {
      if (values.id) {
        updateUserMutation.mutate(values, {
          onSuccess: () => {
            setUserModalProps((prev) => ({ ...prev, show: false }))
            queryClient.invalidateQueries({ queryKey: ['users'] });
          },
          onSettled: () => callback(false)
        });
      } else {
        addUserMutation.mutate(values, {
          onSuccess: () => {
            setUserModalProps((prev) => ({ ...prev, show: false }))
            queryClient.invalidateQueries({ queryKey: ['users'] });
          },
          onSettled: () => callback(false)
        });
      }
    },
    onCancel: () => {
      setUserModalProps((prev) => ({ ...prev, show: false }));
    },
  });

  const columns: ColumnsType<User> = [
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
      title: t('token.user.uniqueName'),
      key: 'uniqueName',
      dataIndex: 'uniqueName',
      align: 'center',
      width: 120,
      render: (text) => (
        <Typography.Text style={{ maxWidth: 120 }} ellipsis={true}>
          {text}
        </Typography.Text>
      )
    },
    {
      title: t('token.user.password'),
      key: 'password',
      dataIndex: 'password',
      align: 'center',
      width: 120,
      render: (text) => (
        <CopyToClipboardInput text={text}/>
      )
    },
    {
      title: t('token.user.enable'),
      key: 'enable',
      dataIndex: 'enable',
      align: 'center',
      render: (status) => {
        if (status === 0) {
          return <Tooltip title={t('token.disable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (status === 1) {
          return <Tooltip title={t('token.enable')}><CheckCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        }
      },
    },
    {
      title: t('token.user.openai'),
      key: 'openai',
      dataIndex: 'openai',
      align: 'center',
      render: (status) => {
        if (status === 0) {
          return <Tooltip title={t('token.disable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (status === 1) {
          return <Tooltip title={t('token.enable')}><CheckCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        }
      },
    },
    // {
    //   title: t('token.user.openaiToken'),
    //   key: 'openaiToken',
    //   dataIndex: 'openaiToken',
    //   ellipsis: true,
    //   align: 'center',
    //   width: 50,
    //   render: (text) => (
    //     <Typography.Text style={{ maxWidth: 50 }} ellipsis={true}>
    //       {text}
    //     </Typography.Text>
    //   )
    // },
    {
      title: t('token.user.claude'),
      key: 'claude',
      dataIndex: 'claude',
      align: 'center',
      render: (status) => {
        if (status === 0) {
          return <Tooltip title={t('token.disable')}><CloseCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (status === 1) {
          return <Tooltip title={t('token.enable')}><CheckCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        }
      },
    },
    {
      title: t("token.expirationTime"),
      key: 'expirationTime',
      dataIndex: 'expirationTime',
      align: 'center',
      width: 200,
      render: (text) => formatDateTime(text),
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
          <Button onClick={() => onEdit(record)} icon={<EditOutlined />} type="primary" />
          <Popconfirm title={t('common.deleteConfirm')} okText="Yes" cancelText="No" placement="left" onConfirm={() => {
            setDeleteUserId(record.id);
            deleteUserMutation.mutate(record.id, {
              onSuccess: () => {
                setDeleteUserId(undefined)
                queryClient.invalidateQueries({ queryKey: ['users'] });
              }
            })
          }}>
            <Button icon={<DeleteOutlined />} type="primary" loading={deleteUserId === record.id} danger />
          </Popconfirm>
        </Button.Group>
      ),
    },
  ];

  useEffect(() => {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(visibleColumns));
  }, [visibleColumns]);

  const handleVisibilityChange = (checkedValues: (keyof User | 'operation')[]) => {
    setTempVisibleColumns(checkedValues);
  };

  const applyColumnVisibility = () => {
    setVisibleColumns(tempVisibleColumns);
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(tempVisibleColumns));
    setPopoverVisible(false);
  };

  const columnVisibilityContent = (
    <div style={{ maxWidth: 110 }}>
      <Checkbox.Group
        options={columns.map(col => ({ label: col.title, value: col.key })) as { label: React.ReactNode; value: keyof User | 'operation' }[]}
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
    col.key && visibleColumns.includes(col.key as keyof User | 'operation')
  );

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['users', uniqueName],
    queryFn: () => userService.searchUserList(uniqueName),
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  });

  const onSearchFormReset = () => {
    searchForm.resetFields();
  };

  const handleRefresh = () => {
    refetch();
    message.success(t('common.dataRefreshed'));
  };

  const onCreate = () => {
    setUserModalProps((prev) => ({
      ...prev,
      show: true,
      title: t('token.createNew'),
      formValue: {
        id: undefined,
        uniqueName: '',
        password: '',
        enable: 1,
        openai: 1,
        claude: 1,
        expirationTime: undefined,
      },
    }));
  };

  const onEdit = (record: UserAddReq) => {
    setUserModalProps({
      formValue: record,
      title: t('token.edit'),
      show: true,
      isEdit: true,
      onOk: (values: UserAddReq, callback) => {
        updateUserMutation.mutate(values, {
          onSuccess: () => {
            setUserModalProps((prev) => ({...prev, show: false}));
            queryClient.invalidateQueries({ queryKey: ['users'] });
          },
          onSettled: () => callback(false)
        })
      },
      onCancel: () => {
        setUserModalProps((prev) => ({...prev, show: false}));
      },
    })
  }

  return (
    <Space direction="vertical" size="large" className="w-full">
      <Card>
        <Form form={searchForm}>
          <Row gutter={[16, 16]}>
            <Col span={6} lg={6}>
              <Form.Item<SearchFormFieldType> label={t('token.user.uniqueName')} name="uniqueName" className="!mb-0">
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
      <UserModal {...userModalProps} />
    </Space>
  );
}

type UserModalProps = {
  formValue: UserAddReq;
  title: string;
  show: boolean;
  isEdit: boolean;
  onOk: (values: UserAddReq, setLoading: (loading: boolean) => void) => void;
  onCancel: VoidFunction;
};

function UserModal({title, show, formValue, onOk, onCancel}: UserModalProps) {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation();

  const [showOpenAI, setShowOpenAI] = useState(formValue.openai === 1);
  const [openAITokens, setOpenAITokens] = useState<Array<{ id: number; tokenName: string }>>([]);
  const [loadingOpenaiAccounts, setLoadingOpenaiAccounts] = useState(false);

  const [showClaude, setShowClaude] = useState(formValue.claude === 1);
  const [claudeTokens, setClaudeTokens] = useState<Array<{ id: number; tokenName: string }>>([]);
  const [loadingClaudeAccounts, setLoadingClaudeAccounts] = useState(false);

  useEffect(() => {
    if (show) {
      form.setFieldsValue({
        ...formValue,
        expirationTime: formValue.expirationTime
          ? dayjs(formValue.expirationTime)
          : dayjs('23:59:59', 'HH:mm:ss').add(30, 'day').tz('Asia/Shanghai')
      });
      setShowOpenAI(formValue.openai === 1);
      setShowClaude(formValue.claude === 1);
    } else {
      form.resetFields();
    }
  }, [show, formValue, form]);

  useEffect(() => {
    if (showOpenAI) {
      setLoadingOpenaiAccounts(true);
      tokenService.searchTokenList('')
        .then(tokens => {
          setOpenAITokens(tokens);
          if (formValue.id && formValue.id > 0 && formValue.openaiToken) {
            form.setFieldsValue({ openaiToken: formValue.openaiToken });
          } else if (tokens.length > 0) {
            form.setFieldsValue({ openaiToken: tokens[0].id });
          }
        })
        .finally(() => setLoadingOpenaiAccounts(false));
    }

    if (showClaude) {
      setLoadingClaudeAccounts(true);
      claudeTokenService.searchTokenList('')
        .then(tokens => {
          setClaudeTokens(tokens);
          if (formValue.id && formValue.id > 0 && formValue.claudeToken) {
            form.setFieldsValue({ claudeToken: formValue.claudeToken });
          } else if (tokens.length > 0) {
            form.setFieldsValue({ claudeToken: tokens[0].id });
          }
        })
        .finally(() => setLoadingClaudeAccounts(false));
    }
  }, [showOpenAI, showClaude, form, formValue.openaiToken, formValue.claudeToken, formValue.id]);

  const onModalOk = () => {
    form.validateFields().then((values) => {
      const formattedValues = {
        ...values,
        expirationTime: values.expirationTime
          ? (dayjs.isDayjs(values.expirationTime)
              ? values.expirationTime
              : dayjs(values.expirationTime)
          ).tz('Asia/Shanghai').format('YYYY-MM-DD HH:mm:ss')
          : null,
      };
      setLoading(true);
      onOk(formattedValues, setLoading);
    }).catch(error => {
      console.error('Validation error:', error);
    });
  };

  const handleOpenAIChange = (value: number) => {
    setShowOpenAI(value === 1);
    if (value === 0) {
      form.setFieldsValue({ openaiToken: 0 });
    }
  };

  const handleClaudeChange = (value: number) => {
    setShowClaude(value === 1);
    if (value === 0) {
      form.setFieldsValue({ claudeToken: 0 });
    }
  };

  return (
    <Modal
      title={title}
      open={show}
      onOk={onModalOk}
      onCancel={() => {
        form.resetFields();
        onCancel();
      }}
      okButtonProps={{
        loading: loading,
      }}
      destroyOnClose={true}
    >
      <Form
        form={form}
        layout="vertical"
      >
        <Form.Item<UserAddReq> name="id" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<UserAddReq> label={t("token.user.uniqueName")} name="uniqueName" required>
          <Input autoComplete="off"/>
        </Form.Item>
        <Form.Item<UserAddReq> label={t("token.user.password")} name="password" required>
          <Input autoComplete="off"/>
        </Form.Item>
        <Form.Item<UserAddReq> label={t("token.user.enable")} name="enable" required>
          <Select>
            <Option value={0}>{t("token.disable")}</Option>
            <Option value={1}>{t("token.enable")}</Option>
          </Select>
        </Form.Item>
        <Form.Item<UserAddReq> label={t("token.user.openai")} name="openai" required>
          <Select onChange={handleOpenAIChange}>
            <Option value={0}>{t("token.disable")}</Option>
            <Option value={1}>{t("token.enable")}</Option>
          </Select>
        </Form.Item>
        {showOpenAI && (
          <Form.Item<UserAddReq> label={t("token.user.openaiToken")} name="openaiToken" required>
            <Select loading={loadingOpenaiAccounts}>
              {openAITokens.map(token => (
                <Option key={token.id} value={token.id}>{token.tokenName}</Option>
              ))}
            </Select>
          </Form.Item>
        )}
        <Form.Item<UserAddReq> label={t("token.user.claude")} name="claude" required>
          <Select onChange={handleClaudeChange}>
            <Option value={0}>{t("token.disable")}</Option>
            <Option value={1}>{t("token.enable")}</Option>
          </Select>
        </Form.Item>
        {showClaude && (
          <Form.Item<UserAddReq> label={t("token.user.claudeToken")} name="claudeToken" required>
            <Select loading={loadingClaudeAccounts}>
              {claudeTokens.map(token => (
                <Option key={token.id} value={token.id}>{token.tokenName}</Option>
              ))}
            </Select>
          </Form.Item>
        )}
        <Form.Item label={t('token.expirationTime')} name="expirationTime" required>
          <DatePicker
            style={{ width: '100%' }}
            format="YYYY-MM-DD HH:mm:ss"
            disabledDate={current => current && current < dayjs().endOf('day')}
            showTime={{ defaultValue: dayjs('23:59:59', 'HH:mm:ss') }}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
}