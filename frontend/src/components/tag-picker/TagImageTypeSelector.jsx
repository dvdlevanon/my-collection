import { ToggleButton, ToggleButtonGroup } from '@mui/material';
import React from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';

function TagImageTypeSelector({ tit, onTitChanged }) {
	const titsQuery = useQuery({
		queryKey: ReactQueryUtil.TAG_IMAGE_TYPES_KEY,
		queryFn: Client.getTagImageTypes,
		onSuccess: (tits) => {
			if (!tit && tits.length > 0) {
				onTitChanged(tits[0]);
			}
		},
	});

	return (
		<>
			{titsQuery.isSuccess && (
				<ToggleButtonGroup
					size="small"
					exclusive
					value={tit}
					onChange={(e, newValue) => {
						if (newValue != null) {
							onTitChanged(newValue);
						}
					}}
				>
					{titsQuery.data.map((tit) => {
						return (
							<ToggleButton key={tit.id} value={tit}>
								<img width="20" height="20" src={Client.buildFileUrl(tit.iconUrl)}></img>
							</ToggleButton>
						);
					})}
				</ToggleButtonGroup>
			)}
		</>
	);
}

export default TagImageTypeSelector;
