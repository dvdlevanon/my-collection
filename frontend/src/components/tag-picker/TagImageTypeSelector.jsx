import { useTheme } from '@emotion/react';
import { ToggleButton, ToggleButtonGroup } from '@mui/material';
import React, { useEffect, useState } from 'react';
import Client from '../../utils/client';
import Svg from '../svg/Svg';

function TagImageTypeSelector({ disabled, tits, tit, onTitChanged, onTitClicked }) {
	const theme = useTheme();
	const [clickedId, setClickedId] = useState(-1);
	const [selectedId, setSelectedId] = useState(-1);

	useEffect(() => {
		if (tit) {
			setSelectedId(tit.id);
		}
	}, [tit]);

	return (
		<ToggleButtonGroup
			disabled={disabled}
			size="small"
			exclusive
			value={tit}
			onChange={(e, newValue) => {
				if (newValue != null) {
					onTitChanged(newValue);
					setSelectedId(newValue.id);
					e.stopPropagation();
				}
			}}
			onClick={(e) => {
				if (!onTitClicked) {
					return;
				}

				if (clickedId > -1) {
					setClickedId(-1);
				} else {
					setClickedId(selectedId);
				}
				onTitClicked(clickedId == -1);
			}}
		>
			{tits.map((tit) => {
				return (
					<ToggleButton key={tit.id} value={tit}>
						<Svg
							width={theme.iconSize(1)}
							height={theme.iconSize(1)}
							color={
								clickedId == tit.id && onTitClicked
									? theme.palette.text.light
									: theme.palette.text.primary
							}
							path={Client.buildFileUrl(tit.iconUrl)}
						/>
					</ToggleButton>
				);
			})}
		</ToggleButtonGroup>
	);
}

export default TagImageTypeSelector;
