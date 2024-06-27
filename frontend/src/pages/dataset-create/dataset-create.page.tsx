import { Flex, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';
import styles from './dataset-create.module.scss';
import { DatasetForm } from '@/components/datasets/components/dataset-form.component.tsx';

export function DatasetCreate() {
    const { t } = useTranslation();
    return (
        <Flex vertical className={styles.container}>
            <Typography.Title level={3}>{t('New Dataset')}</Typography.Title>
            <Space direction={'vertical'} size={'large'} className={styles.container}>
                <DatasetForm mode={'create'} />
            </Space>
        </Flex>
    );
}
