import { Box, Link, Stack, Typography } from '@mui/material';
import { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../utils/client';
import ItemsUtil from '../../utils/items-util';
import TimeUtil from '../../utils/time-utils';
import ItemBadges from './ItemBadges';
import ItemCoverIndicator from './ItemCoverIndicator';
import ItemOffests from './ItemOffests';
import ItemTitle from './ItemTitle';

function Item({
	item,
	preferPreview,
	itemWidth,
	itemHeight,
	direction,
	showOffests,
	titleSx,
	withItemTitleMenu,
	itemLinkBuilder,
	onConvertAudio,
	onConvertVideo,
}) {
	let [mouseEnterMillis, setMouseEnterMillis] = useState(0);
	let [showCoverNavigator, setShowCoverNavigator] = useState(false);
	let [showPreview, setShowPreview] = useState(false);
	let [coverNumber, setCoverNumber] = useState(
		item.covers && item.covers.length > 0 ? Math.floor(item.covers.length / 2) : 0
	);

	const shouldShowPreview = () => {
		return preferPreview && ItemsUtil.hasPreview(item);
	};

	const mouseMove = (e) => {
		if (shouldShowPreview() || !item.covers) {
			return;
		}

		let nowMillis = Math.floor(Date.now());

		if (mouseEnterMillis === 0 || mouseEnterMillis > nowMillis - 250) {
			return;
		}

		let bounds = e.currentTarget.getBoundingClientRect();
		let x = e.clientX - bounds.left;
		setShowCoverNavigator(true);
		setCoverNumber(Math.floor(x / (bounds.width / item.covers.length)));
	};

	const mouseLeave = (e) => {
		if (shouldShowPreview()) {
			setShowPreview(false);
		} else {
			setMouseEnterMillis(0);
			setCoverNumber(item.covers && item.covers.length > 0 ? Math.floor(item.covers.length / 2) : 0);
			setShowCoverNavigator(false);
		}
	};

	const mouseEnter = (e) => {
		if (shouldShowPreview()) {
			setShowPreview(true);
		} else {
			setMouseEnterMillis(Math.floor(Date.now()));
		}
	};

	const getTitleComponent = () => {
		return (
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'column',
					padding: '10px',
					gap: '10px',
				}}
			>
				<ItemTitle
					item={item}
					variant="caption"
					withTooltip={true}
					withMenu={withItemTitleMenu}
					preventDefault={direction === 'column' ? true : false}
					onTagAdded={() => console.log('unsupported adding tag')}
					onTitleChanged={() => console.log('unsupported changing title')}
					sx={{
						whiteSpace: 'nowrap',
						overflow: 'hidden',
						textAlign: 'center',
						...titleSx,
					}}
				/>
				{showOffests && item.main_item && <ItemOffests item={item} />}
			</Box>
		);
	};

	return (
		<Stack
			flexDirection={direction}
			sx={{
				maxWidth: direction === 'column' ? itemWidth : 'unset',
			}}
		>
			<Link
				component={RouterLink}
				reloadDocument
				to={itemLinkBuilder(item)}
				sx={{
					display: 'flex',
					position: 'relative',
					flexDirection: 'column',
					width: itemWidth,
					height: itemHeight,
				}}
				onMouseLeave={(e) => mouseLeave(e)}
				onMouseMove={(e) => mouseMove(e)}
				onMouseEnter={(e) => mouseEnter(e)}
			>
				<ItemBadges item={item} onConvertAudio={onConvertAudio} onConvertVideo={onConvertVideo} />
				<Box position="relative">
					<Box
						component="img"
						src={ItemsUtil.getCover(item, coverNumber)}
						alt={item.title}
						loading="lazy"
						sx={{
							width: itemWidth,
							height: itemHeight,
							objectFit: 'contain',
							cursor: 'pointer',
							borderRadius: '10px',
							boxShadow: item.main_cover_second ? '0 2px 4px 0 rgba(13, 147, 82, 0.6)' : 'none',
						}}
					/>
					{!ItemsUtil.isHighlight(item) && (
						<Typography
							variant="caption"
							sx={{
								position: 'absolute',
								backgroundColor: 'black',
								color: 'white',
								right: '3px',
								bottom: '10px',
								borderRadius: '5px',
								opacity: 0.9,
								padding: '0px 2px',
								zIndex: 100,
							}}
						>
							{TimeUtil.formatDuration(item.duration_seconds)}
						</Typography>
					)}
				</Box>
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
									isHighlighted={coverNumber === index}
								/>
							);
						})}
					</Stack>
				)}
				{shouldShowPreview() && showPreview && (
					<Box
						flexGrow={1}
						padding="10px"
						sx={{
							position: 'absolute',
							padding: '0px',
							width: itemWidth,
							height: itemHeight,
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
							onLoadedMetadata={(e) => {
								if (item.preview_mode === ItemsUtil.PREVIEW_FROM_START_POSITION) {
									e.target.currentTime = item.start_position;
								}
							}}
							onTimeUpdate={(e) => {
								if (
									item.preview_mode === ItemsUtil.PREVIEW_FROM_START_POSITION &&
									e.target.currentTime > item.end_position
								) {
									e.target.currentTime = item.start_position;
								}
							}}
							sx={{
								borderRadius: '10px',
							}}
						>
							<source src={Client.buildFileUrl(ItemsUtil.getPreview(item))} />
						</Box>
					</Box>
				)}
			</Link>
			{!ItemsUtil.isHighlight(item) && direction === 'column' && getTitleComponent()}
			{!ItemsUtil.isHighlight(item) && direction !== 'column' && (
				<Link component={RouterLink} reloadDocument to={itemLinkBuilder(item)} color="bright.text">
					{getTitleComponent()}
				</Link>
			)}
		</Stack>
	);
}

export default Item;
