import { Box, Tooltip } from '@mui/material';
import { useState } from 'react';
import { Link, Link as RouterLink } from 'react-router-dom';
import seedrandom from 'seedrandom';
import Client from '../../utils/client';
import GalleryUrlParams from '../../utils/gallery-url-params';
import ManageTagImageDialog from '../dialogs/ManageTagImageDialog';
import TagContextMenu from '../tag-context-menu/TagContextMenu';
import Thumbnail from '../thumbnail/Thumbnail';

function TagThumbnail({ tag, isLink, onTagClicked, onTagRemoved, withRemoveOption }) {
	const [tagMenuProps, setTagMenuProps] = useState(null);
	const [manageTagImageOpened, setManageTagImageOpened] = useState(false);
	const [autoThumbnailMode, setAutoThumbnailMode] = useState(false);

	const getRandomImage = (images) => {
		let epochDay = Math.floor(Date.now() / 1000 / 60 / 60 / 24);
		let rand = seedrandom(epochDay + tag.id);
		let randomIndex = Math.floor(rand() * images.length);
		return images[randomIndex];
	};

	const getThumbnailCompnent = () => {
		if (!tag.images) {
			return <Thumbnail title={tag.title} />;
		}

		let imagesWithThumbnails = tag.images.filter((image) => {
			return image.thumbnail_rect && image.thumbnail_rect.height != 0;
		});

		if (imagesWithThumbnails.length == 0) {
			return <Thumbnail title={tag.title} />;
		}

		let image = getRandomImage(imagesWithThumbnails);
		return <Thumbnail crop={image.thumbnail_rect} imageUrl={Client.buildFileUrl(image.url)} />;
	};

	const getThumbnailCompnentWrapper = () => {
		return (
			<Tooltip title={tag.title}>
				<Box
					onContextMenu={(e) => {
						e.preventDefault();
						setTagMenuProps({
							anchor: e.target,
							left: e.clientX,
							top: e.clientY,
						});
					}}
					onClick={() => {
						if (onTagClicked) {
							onTagClicked(tag);
						}
					}}
					sx={{
						cursor: 'pointer',
					}}
				>
					{getThumbnailCompnent()}
				</Box>
			</Tooltip>
		);
	};

	return (
		<>
			{isLink && (
				<Link target="_blank" component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
					{getThumbnailCompnentWrapper()}
				</Link>
			)}
			{!isLink && getThumbnailCompnentWrapper()}
			{tagMenuProps != null && (
				<TagContextMenu
					tag={tag}
					menuAnchorEl={tagMenuProps.anchor}
					menuPosition={{ top: tagMenuProps.top, left: tagMenuProps.left }}
					onClose={() => setTagMenuProps(null)}
					onManageAttributesClicked={null}
					withManageAttributesClicked={false}
					onRemoveTagClicked={() => onTagRemoved(tag)}
					withRemoveOption={withRemoveOption}
					onManageImageClicked={() => {
						setManageTagImageOpened(true);
						setAutoThumbnailMode(false);
					}}
					onEditThumbnail={() => {
						setManageTagImageOpened(true);
						setAutoThumbnailMode(true);
					}}
				/>
			)}
			{manageTagImageOpened && (
				<ManageTagImageDialog
					tag={tag}
					autoThumbnailMode={autoThumbnailMode}
					onClose={() => setManageTagImageOpened(false)}
				/>
			)}
		</>
	);
}

export default TagThumbnail;
