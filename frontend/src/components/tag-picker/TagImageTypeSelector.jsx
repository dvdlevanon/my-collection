import { ToggleButton, ToggleButtonGroup } from '@mui/material';
import React from 'react';
import Client from '../../utils/client';

function TagImageTypeSelector({ tits, tit, onTitChanged }) {
	return (
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
			{tits.map((tit) => {
				return (
					<ToggleButton key={tit.id} value={tit}>
						<img width="20" height="20" src={Client.buildFileUrl(tit.iconUrl)}></img>
					</ToggleButton>
				);
			})}
		</ToggleButtonGroup>
	);
}

export default TagImageTypeSelector;
