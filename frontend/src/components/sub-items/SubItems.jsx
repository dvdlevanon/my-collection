import { Stack } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import SubItem from './SubItem';

function SubItems({ item }) {
	const mainItemQuery = useQuery({
		queryKey: ReactQueryUtil.itemKey(item.main_item || item.id),
		queryFn: () => {
			if (item.main_item) {
				return Client.getItem(item.main_item);
			} else {
				return item;
			}
		},
	});

	const sortedSubItems = () => {
		return mainItemQuery.data.sub_items.sort((item1, item2) => {
			return item1.start_position > item2.start_position
				? 1
				: item1.start_position < item2.start_position
				? -1
				: 0;
		});
	};

	return (
		<>
			{mainItemQuery.isSuccess && (
				<Stack flexDirection="column" gap="10px">
					{sortedSubItems().map((subItem) => {
						return (
							<SubItem
								key={subItem.id}
								item={subItem}
								itemWidth={200}
								highlighted={subItem.id == item.id}
							/>
						);
					})}
				</Stack>
			)}
		</>
	);
}

export default SubItems;
