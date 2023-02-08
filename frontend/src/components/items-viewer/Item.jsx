import { Box, Link, Stack, Tooltip, Typography } from '@mui/material';
import { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../network/client';
import ItemCoverIndicator from './ItemCoverIndicator';

function Item({ item, preferPreview }) {
	let [mouseEnterMillis, setMouseEnterMillis] = useState(0);
	let [showCoverNavigator, setShowCoverNavigator] = useState(false);
	let [showPreview, setShowPreview] = useState(false);
	let [coverNumber, setCoverNumber] = useState(
		item.covers && item.covers.length > 0 ? Math.floor(item.covers.length / 2) : 0
	);

	const getCover = () => {
		if (item.covers && item.covers.length > 0 && item.covers[coverNumber]) {
			return Client.buildFileUrl(item.covers[coverNumber].url);
		} else {
			return 'empty';
		}
	};

	const previewMode = () => {
		return preferPreview && item.previewUrl != '';
	};

	const mouseMove = (e) => {
		if (previewMode() || !item.covers) {
			return;
		}

		let nowMillis = Math.floor(Date.now());

		if (mouseEnterMillis == 0 || mouseEnterMillis > nowMillis - 250) {
			return;
		}

		let bounds = e.currentTarget.getBoundingClientRect();
		let x = e.clientX - bounds.left;
		setShowCoverNavigator(true);
		setCoverNumber(Math.floor(x / (bounds.width / item.covers.length)));
	};

	const mouseLeave = (e) => {
		if (previewMode()) {
			setShowPreview(false);
		} else {
			setMouseEnterMillis(0);
			setCoverNumber(item.covers && item.covers.length > 0 ? Math.floor(item.covers.length / 2) : 0);
			setShowCoverNavigator(false);
		}
	};

	const mouseEnter = (e) => {
		if (previewMode()) {
			setShowPreview(true);
		} else {
			setMouseEnterMillis(Math.floor(Date.now()));
		}
	};

	const getFormattedDuration = () => {
		if (!item.duration_seconds) {
			return '00:00';
		}

		if (item.duration_seconds < 60 * 60) {
			return new Date(item.duration_seconds * 1000).toISOString().slice(14, 19);
		} else {
			return new Date(item.duration_seconds * 1000).toISOString().slice(11, 19);
		}
	};

	return (
		<Link
			component={RouterLink}
			sx={{
				display: 'flex',
				position: 'relative',
				flexDirection: 'column',
			}}
			to={'/spa/item/' + item.id}
			onMouseLeave={(e) => mouseLeave(e)}
			onMouseMove={(e) => mouseMove(e)}
			onMouseEnter={(e) => mouseEnter(e)}
		>
			<Box
				component="img"
				src={getCover()}
				alt={item.title}
				loading="lazy"
				sx={{
					width: '550px',
					height: '300px',
					objectFit: 'contain',
					cursor: 'pointer',
				}}
			/>
			{showCoverNavigator && item.covers && item.covers.length > 1 && (
				<Stack
					flexDirection="row"
					sx={{
						bottom: '0px',
						left: '0px',
						gap: '2px',
						position: 'absolute',
						width: '100%',
					}}
				>
					{item.covers.map((cover, index) => {
						return (
							<ItemCoverIndicator
								key={cover.id}
								item={item}
								cover={cover}
								isHighlighted={coverNumber == index}
							/>
						);
					})}
				</Stack>
			)}
			{previewMode() && showPreview && (
				<Box
					flexGrow={1}
					padding="10px"
					sx={{
						position: 'absolute',
						padding: '0px',
						width: '550px',
						height: '300px',
						objectFit: 'contain',
						cursor: 'pointer',
					}}
				>
					<Box
						component="video"
						height="100%"
						width="100%"
						playsInline
						muted
						autoPlay={true}
						loop={true}
						controls={false}
					>
						<source src={Client.buildFileUrl(item.previewUrl)} />
					</Box>
				</Box>
			)}
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
					gap: '10px',
					alignItems: 'center',
					justifyContent: 'center',
					height: '50px',
				}}
			>
				<Typography
					variant="caption"
					sx={{
						padding: '0px 3px',
						borderWidth: '1px',
						borderColor: 'bright.main',
						borderStyle: 'solid',
						borderRadius: '3px',
						color: 'bright.main',
						verticalAlign: 'middle',
						margin: '10px',
					}}
				>
					{getFormattedDuration()}
				</Typography>
				<Tooltip title={item.title} arrow followCursor>
					<Typography
						variant="caption"
						sx={{
							whiteSpace: 'nowrap',
							overflow: 'hidden',
							textOverflow: 'ellipsis',
							cursor: 'pointer',
							maxWidth: '450px',
							textAlign: 'center',
							padding: '5px',
							color: 'primary.light',
							flexGrow: 1,
						}}
					>
						{item.title}
					</Typography>
				</Tooltip>
			</Box>
		</Link>
	);
}

export default Item;
