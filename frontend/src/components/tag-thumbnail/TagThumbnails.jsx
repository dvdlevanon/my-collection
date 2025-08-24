import { useTheme } from '@emotion/react';
import { Stack } from '@mui/material';
import ScrollContainer from 'react-indiana-drag-scroll';
import 'react-indiana-drag-scroll/dist/style.css';
import TagThumbnail from './TagThumbnail';

function TagThumbnails({ tags, carouselMode, onTagClicked, onTagRemoved, withRemoveOption }) {
	const theme = useTheme();

	const getTagComponents = () => {
		return tags.map((tag) => {
			return (
				<TagThumbnail
					key={tag.id}
					tag={tag}
					onTagClicked={(tag) => {
						if (onTagClicked) {
							onTagClicked(tag);
						}
					}}
					onTagRemoved={onTagRemoved}
					withRemoveOption={withRemoveOption}
					isLink={onTagClicked == null}
				/>
			);
		});
	};

	return (
		<>
			{!carouselMode && (
				<Stack flexDirection="row" gap={theme.spacing(1)} alignItems="center">
					{getTagComponents()}
				</Stack>
			)}
			{carouselMode && (
				<ScrollContainer vertical="false">
					<Stack height={theme.iconSize(3)} width="100%" flexDirection="row" gap={theme.spacing(1)}>
						{getTagComponents()}
					</Stack>
				</ScrollContainer>
			)}
		</>
	);
}

export default TagThumbnails;
