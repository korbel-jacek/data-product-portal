import { Flex, Table, type TableColumnsType } from 'antd';
import { useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

import {
    useGetDataOutputByIdQuery,
    useRemoveDatasetFromDataOutputMutation,
} from '@/store/features/data-outputs/data-outputs-api-slice.ts';
import { dispatchMessage } from '@/store/features/feedback/utils/dispatch-feedback.ts';
import type { DataOutputDatasetLink } from '@/types/data-output';

import styles from './dataset-table.module.scss';
import { getDataOutputDatasetsColumns } from './dataset-table-columns.tsx';

type Props = {
    isCurrentDataOutputOwner: boolean;
    dataOutputId: string;
    datasets: DataOutputDatasetLink[];
};

export function DatasetTable({ isCurrentDataOutputOwner, dataOutputId, datasets }: Props) {
    const { t } = useTranslation();
    const { data: dataOutput, isLoading: isLoadingDataOutput } = useGetDataOutputByIdQuery(dataOutputId);
    const [removeDatasetFromDataOutput, { isLoading: isRemovingDatasetFromDataOutput }] =
        useRemoveDatasetFromDataOutputMutation();

    const handleRemoveDatasetFromDataOutput = useCallback(
        async (datasetId: string, name: string) => {
            try {
                await removeDatasetFromDataOutput({ datasetId, dataOutputId: dataOutputId }).unwrap();
                dispatchMessage({
                    content: t('Dataset {{name}} has been removed from data output', { name }),
                    type: 'success',
                });
            } catch (error) {
                console.error('Failed to remove dataset from data output', error);
            }
        },
        [dataOutputId, removeDatasetFromDataOutput, t],
    );

    const handleCancelDatasetLinkRequest = useCallback(
        async (datasetId: string, name: string) => {
            try {
                await removeDatasetFromDataOutput({ datasetId, dataOutputId: dataOutputId }).unwrap();
                dispatchMessage({
                    content: t('Request to link dataset {{name}} has been cancelled', { name }),
                    type: 'success',
                });
            } catch (error) {
                console.error('Failed to cancel dataset link request', error);
            }
        },
        [dataOutputId, removeDatasetFromDataOutput, t],
    );

    const columns: TableColumnsType<DataOutputDatasetLink> = useMemo(() => {
        return getDataOutputDatasetsColumns({
            onRemoveDataOutputDatasetLink: handleRemoveDatasetFromDataOutput,
            onCancelDataOutputDatasetLinkRequest: handleCancelDatasetLinkRequest,
            t,
            isDisabled: !isCurrentDataOutputOwner,
            isLoading: isRemovingDatasetFromDataOutput,
        });
    }, [
        handleRemoveDatasetFromDataOutput,
        handleCancelDatasetLinkRequest,
        t,
        isCurrentDataOutputOwner,
        isRemovingDatasetFromDataOutput,
    ]);

    if (!dataOutput) return null;

    return (
        <Flex className={styles.datasetListContainer}>
            <Table<DataOutputDatasetLink>
                loading={isLoadingDataOutput}
                className={styles.datasetListTable}
                columns={columns}
                dataSource={datasets}
                rowKey={({ id }) => id}
                pagination={false}
                rowClassName={styles.tableRow}
                size={'small'}
            />
        </Flex>
    );
}
