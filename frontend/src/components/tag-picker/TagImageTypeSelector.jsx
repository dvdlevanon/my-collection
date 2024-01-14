import { useTheme } from '@emotion/react';
import { ToggleButton, ToggleButtonGroup } from '@mui/material';
import React from 'react';
import Client from '../../utils/client';
import Svg from '../svg/Svg';

function TagImageTypeSelector({ disabled, tits, tit, onTitChanged }) {
	const theme = useTheme();

	return (
		<ToggleButtonGroup
			disabled={disabled}
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
						<Svg
							width={theme.iconSize(1)}
							height={theme.iconSize(1)}
							color={theme.palette.text.primary}
							path={Client.buildFileUrl(tit.iconUrl)}
						/>
					</ToggleButton>
				);
			})}
		</ToggleButtonGroup>
	);
}

export default TagImageTypeSelector;
