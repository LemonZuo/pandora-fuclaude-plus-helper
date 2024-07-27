import { useState, useEffect } from 'react';
import {
  Badge,
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
  Spin,
  Tooltip,
  Typography,
  Checkbox, Popover, CheckboxOptionType, message
} from 'antd';
import Table, { ColumnsType } from 'antd/es/table';
import {
  CheckCircleOutlined,
  DeleteOutlined,
  EditOutlined,
  FundOutlined,
  MinusCircleOutlined,
  QuestionCircleOutlined,
  ReloadOutlined,
  ShareAltOutlined
} from "@ant-design/icons";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import dayjs from "dayjs";
import 'dayjs/locale/zh-cn';
import customParseFormat from 'dayjs/plugin/customParseFormat';
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

import { OpenaiAccount, OpenaiToken } from '#/entity.ts';
import tokenService, { OpenaiTokenAddReq } from "@/api/services/tokenService.ts";
import accountService from "@/api/services/accountService.ts";
import {
  useAddTokenMutation,
  useDeleteTokenMutation,
  useRefreshTokenMutation,
  useUpdateTokenMutation
} from "@/store/tokenStore.ts";
import { useAddAccountMutation } from "@/store/accountStore.ts";
import CopyToClipboardInput from "@/pages/components/copy";
import formatDateTime from "@/pages/components/util";
import Chart from "@/components/chart/chart.tsx";
import useChart from "@/components/chart/useChart.ts";

dayjs.locale('zh-cn');
dayjs.extend(utc);
dayjs.extend(timezone);
dayjs.extend(customParseFormat);

type SearchFormFieldType = Pick<OpenaiToken, 'tokenName'>;

const { Option } = Select;

const LOCAL_STORAGE_KEY = 'openai_token_page_visible_columns';

export default function TokenPage() {
  const queryClient = useQueryClient();
  const [searchForm] = Form.useForm();
  const { t } = useTranslation();

  const addTokenMutation = useAddTokenMutation();
  const updateTokenMutation = useUpdateTokenMutation();
  const deleteTokenMutation = useDeleteTokenMutation();
  const refreshTokenMutation = useRefreshTokenMutation();
  const addAccountMutation = useAddAccountMutation();

  const navigate = useNavigate();

  const [deleteTokenId, setDeleteTokenId] = useState<number | undefined>(-1);
  const [refreshTokenId, setRefreshTokenId] = useState<number | undefined>(-1);

  const [visibleColumns, setVisibleColumns] = useState<(keyof OpenaiToken | 'operation' | 'share')[]>(() => {
    const storedColumns = localStorage.getItem(LOCAL_STORAGE_KEY);
    return storedColumns
      ? JSON.parse(storedColumns)
      : ['id', 'tokenName', 'plusSubscription', 'refreshToken', 'accessToken',
        'expireAt', 'createTime', 'updateTime', 'share', 'operation'];
  });
  const [tempVisibleColumns, setTempVisibleColumns] = useState<(keyof OpenaiToken | 'operation' | 'share')[]>(visibleColumns);
  const [popoverVisible, setPopoverVisible] = useState(false);

  const searchTokenName = Form.useWatch('tokenName', searchForm);

  const { data, isLoading, refetch } = useQuery({
    queryKey: ['openaiTokens', searchTokenName],
    queryFn: () => tokenService.searchTokenList(searchTokenName),
    refetchOnMount: true,
    refetchOnWindowFocus: true,
  });

  const [TokenModalProps, setTokenModalProps] = useState<TokenModalProps>({
    formValue: {
      tokenName: '',
      refreshToken: '',
    },
    title: 'New',
    show: false,
    onOk: (values: OpenaiTokenAddReq, callback) => {
      if (values.id) {
        updateTokenMutation.mutate(values, {
          onSuccess: () => {
            setTokenModalProps((prev) => ({ ...prev, show: false }))
            queryClient.invalidateQueries({ queryKey: ['openaiTokens'] });
          },
          onSettled: () => callback(false)
        });
      } else {
        addTokenMutation.mutate(values, {
          onSuccess: () => {
            setTokenModalProps((prev) => ({ ...prev, show: false }))
            queryClient.invalidateQueries({ queryKey: ['openaiTokens'] });
          },
          onSettled: () => callback(false)
        });
      }
    },
    onCancel: () => {
      setTokenModalProps((prev) => ({ ...prev, show: false }));
    },
  });

  const [shareModalProps, setAccountModalProps] = useState<AccountModalProps>({
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
    title: 'New',
    show: false,
    isEdit: false,
    onOk: (values: OpenaiAccount, callback) => {
      values.gpt35Limit = parseInt(values.gpt35Limit as any);
      values.gpt4Limit = parseInt(values.gpt4Limit as any);
      callback(true);
      addAccountMutation.mutate(values, {
        onSuccess: () => {
          setAccountModalProps((prev) => ({ ...prev, show: false }))
          queryClient.invalidateQueries({ queryKey: ['openaiTokens'] });
        },
        onSettled: () => callback(false)
      });
    },
    onCancel: () => {
      setAccountModalProps((prev) => ({ ...prev, show: false }));
    },
  });

  const [shareInfoModalProps, setAccountInfoModalProps] = useState<AccountInfoModalProps>({
    tokenId: -1,
    show: false,
    onOk: () => {
      setAccountInfoModalProps((prev) => ({ ...prev, show: false }));
    },
  });

  const columns: ColumnsType<OpenaiToken> = [
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
      render: (text) => (
        <Typography.Text style={{ maxWidth: 200 }} ellipsis={true}>
          {text}
        </Typography.Text>
      )
    },
    {
      title: t('token.plusSubscription'),
      key: 'plusSubscription',
      dataIndex: 'plusSubscription',
      align: 'center',
      render: (subscription) => {
        if (subscription === 1) {
          return <Tooltip title={t('token.subscriptionUnknown')}><QuestionCircleOutlined style={{ color: 'gray' }} /></Tooltip>;
        } else if (subscription === 2) {
          return <Tooltip title={t('token.unsubscribed')}><MinusCircleOutlined style={{ color: 'red' }} /></Tooltip>;
        } else if (subscription === 3) {
          return <Tooltip title={t('token.subscribed')}><CheckCircleOutlined style={{ color: 'green' }} /></Tooltip>;
        }
      },
    },
    {
      title: t('token.refreshToken'),
      key: 'refreshToken',
      dataIndex: 'refreshToken',
      align: 'center',
      render: (text) => (
        <CopyToClipboardInput text={text} />
      ),
    },
    {
      title: t('token.accessToken'),
      key: 'accessToken',
      dataIndex: 'accessToken',
      align: 'center',
      render: (text) => (
        <CopyToClipboardInput text={text} />
      ),
    },
    {
      title: t("token.expireAt"),
      key: 'expireAt',
      dataIndex: 'expireAt',
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
      title: t('token.share'),
      key: 'share',
      align: 'center',
      render: (_, record) => (
        <Button.Group>
          <Badge style={{ zIndex: 9 }}>
            <Button icon={<ShareAltOutlined />} onClick={() => navigate({
              pathname: '/token/openai-account',
              search: `?tokenId=${record.id}`,
            })}>
              {t('token.shareList')}
            </Button>
          </Badge>
          {/*<Button icon={<PlusOutlined />} onClick={() => onAccountAdd(record)} />*/}
          <Button icon={<FundOutlined />} onClick={() => onAccountInfo(record)} />
        </Button.Group>
      ),
    },
    {
      title: t('token.action'),
      key: 'operation',
      align: 'center',
      render: (_, record) => (
        <Button.Group>
          <Popconfirm title={t('common.refreshConfirm')} okText={t('common.yes')} cancelText={t('common.no')} placement="left" onConfirm={() => {
            setRefreshTokenId(record.id);
            refreshTokenMutation.mutate(record.id, {
              onSettled: () => setRefreshTokenId(undefined),
            })
          }}>
            <Button key={record.id} icon={<ReloadOutlined />} type="primary" loading={refreshTokenId === record.id} style={{ backgroundColor: '#007bff', borderColor: '#007bff', color: 'white' }}>{t('common.refresh')}</Button>
          </Popconfirm>
          <Button onClick={() => onEdit(record)} icon={<EditOutlined />} type="primary" />
          <Popconfirm title={t('common.deleteConfirm')} okText={t('common.yes')} cancelText={t('common.no')} placement="left" onConfirm={() => {
            setDeleteTokenId(record.id);
            deleteTokenMutation.mutate(record.id, {
              onSuccess: () => {
                setDeleteTokenId(undefined)
                queryClient.invalidateQueries({ queryKey: ['openaiTokens'] });
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

  const handleVisibilityChange = (checkedValues: (keyof OpenaiToken | 'operation' | 'share')[]) => {
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
        options={columns.map(col => ({ label: col.title, value: col.key })) as CheckboxOptionType<keyof OpenaiToken | "operation" | "share">[]}
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
    col.key && visibleColumns.includes(col.key as keyof OpenaiToken | 'operation' | 'share')
  );

  const onSearchFormReset = () => {
    searchForm.resetFields();
  };

  const handleRefresh = () => {
    refetch();
    message.success(t('common.dataRefreshed'));
  };

  const onCreate = () => {
    setTokenModalProps((prev) => ({
      ...prev,
      show: true,
      title: t('token.createNew'),
      formValue: {
        id: undefined,
        tokenName: '',
        refreshToken: '',
      },
    }));
  };

  const onAccountInfo = (record: OpenaiToken) => {
    setAccountInfoModalProps((prev) => ({
      ...prev,
      show: true,
      isEdit: false,
      tokenId: record.id,
    }));
  }

  const onEdit = (record: OpenaiToken) => {
    setTokenModalProps((prev) => ({
      ...prev,
      show: true,
      title: t('token.edit'),
      formValue: {
        id: record.id,
        tokenName: record.tokenName,
        refreshToken: record.refreshToken,
      },
    }));
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
      <TokenModal {...TokenModalProps} />
      <AccountModal {...shareModalProps} />
      <AccountInfoModal {...shareInfoModalProps} />
    </Space>
  );
}

export type AccountModalProps = {
  formValue: OpenaiAccount;
  title: string;
  show: boolean;
  isEdit: boolean;
  onOk: (values: OpenaiAccount, callback: any) => void;
  onCancel: VoidFunction;
}

export const AccountModal = ({ title, show, isEdit, formValue, onOk, onCancel }: AccountModalProps) => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const { t } = useTranslation()

  useEffect(() => {
    if (show) {
      const expirationTimeWithZone = formValue.expirationTime
        ? dayjs(formValue.expirationTime, 'YYYY-MM-DD HH:mm:ss').tz('Asia/Shanghai')
        : dayjs().add(30, 'day').tz('Asia/Shanghai');

      form.setFieldsValue({
        ...formValue,
        expirationTime: expirationTimeWithZone,
      });
    } else {
      form.resetFields();
    }
  }, [formValue, show, form]);

  const onModalOk = () => {
    form.validateFields().then((values) => {
      const formattedValues = {
        ...values,
        expirationTime: values.expirationTime ? dayjs(values.expirationTime).tz('Asia/Shanghai').format('YYYY-MM-DD HH:mm:ss') : null,
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
        <Form.Item<OpenaiAccount> name="id" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<OpenaiAccount> name="userId" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<OpenaiAccount> name="tokenId" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<OpenaiAccount> name="status" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<OpenaiAccount> label="OpenaiAccount" name="account" required>
          <Input readOnly={isEdit} disabled={isEdit} autoComplete="off"/>
        </Form.Item>
        <Form.Item label={t('token.expirationTime')} name="expirationTime" required>
          <DatePicker
            style={{ width: '100%' }}
            format="YYYY-MM-DD HH:mm:ss"
            disabledDate={current => current && current < dayjs().endOf('day')}
            showTime={{ defaultValue: dayjs('00:00:00', 'HH:mm:ss') }}
            disabled={isEdit}
          />
        </Form.Item>
        <Form.Item<OpenaiAccount> label={t('token.gpt35Limit')} name="gpt35Limit" required>
          <Input/>
        </Form.Item>
        <Form.Item<OpenaiAccount> label={t('token.gpt4Limit')} name="gpt4Limit" required>
          <Input/>
        </Form.Item>
        <Form.Item<OpenaiAccount> label={t('token.showConversations')} name="showConversations" initialValue={0} required>
          <Select allowClear>
            <Option value={1}>{t('common.yes')}</Option>
            <Option value={0}>{t('common.no')}</Option>
          </Select>
        </Form.Item>
        <Form.Item<OpenaiAccount> label={t('token.temporaryChat')} name="temporaryChat" initialValue={0} required>
          <Select allowClear>
            <Option value={1}>{t('common.yes')}</Option>
            <Option value={0}>{t('common.no')}</Option>
          </Select>
        </Form.Item>
      </Form>
    </Modal>
  );
}

type TokenModalProps = {
  formValue: OpenaiTokenAddReq;
  title: string;
  show: boolean;
  onOk: (values: OpenaiTokenAddReq, setLoading: any) => void;
  onCancel: VoidFunction;
};

function TokenModal({title, show, formValue, onOk, onCancel}: TokenModalProps) {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const {t} = useTranslation()

  useEffect(() => {
    if (show) {
      form.setFieldsValue(formValue);
    } else {
      form.resetFields();
    }
  }, [show, formValue, form]);

  const onModalOk = () => {
    form.validateFields().then((values) => {
      setLoading(true)
      onOk(values, setLoading);
    });
  }

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
        <Form.Item<OpenaiTokenAddReq> name="id" hidden>
          <Input/>
        </Form.Item>
        <Form.Item<OpenaiTokenAddReq> label={t("token.tokenName")} name="tokenName" required>
          <Input autoComplete="off"/>
        </Form.Item>
        <Form.Item<OpenaiTokenAddReq> label={t("token.refreshToken")} name="refreshToken" required>
          <Input autoComplete="off"/>
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
    queryKey: ['openaiTokenStatistic', tokenId],
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
