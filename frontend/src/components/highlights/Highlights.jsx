import { Stack } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import Highlight from './Highlight';

function Highlights({ item }) {
	const mainItemQuery = useQuery({
		queryKey: ReactQueryUtil.itemKey(item.highlight_parent_id || item.id),
		queryFn: () => {
			if (item.highlight_parent_id) {
				return Client.getItem(item.highlight_parent_id);
			} else {
				return item;
			}
		},
	});

	return (
		<>
			{mainItemQuery.isSuccess && (
				<Stack flexDirection="column" gap="10px">
					{mainItemQuery.data.id !== item.id && (
						<Highlight
							item={mainItemQuery.data}
							itemWidth={200}
							highlighted={mainItemQuery.data == item.id}
						/>
					)}
					{mainItemQuery.data.highlights.map((highlight) => {
						return (
							<Highlight
								key={highlight.id}
								item={highlight}
								itemWidth={200}
								highlighted={highlight.id == item.id}
							/>
						);
					})}
				</Stack>
			)}
		</>
	);
}

export default Highlights;
