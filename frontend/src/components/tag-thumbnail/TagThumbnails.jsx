import { useTheme } from '@emotion/react';
import { Box, Stack } from '@mui/material';
import { useState } from 'react';
import Carousel from 'react-multi-carousel';
import 'react-multi-carousel/lib/styles.css';
import TagThumbnail from './TagThumbnail';

function TagThumbnails({ tags, carouselMode, onTagClicked, onTagRemoved, withRemoveOption }) {
	const [duringDrag, setDuringDrag] = useState(false);
	const theme = useTheme();

	const responsive = {
		superLargeDesktop: {
			breakpoint: { max: 4000, min: 3000 },
			items: 2560 / (theme.iconBaseSize * 3 + theme.baseSpacing * 2),
		},
		desktop: {
			breakpoint: { max: 3000, min: 1024 },
			items: 2560 / (theme.iconBaseSize * 3 + theme.baseSpacing * 2),
		},
		tablet: {
			breakpoint: { max: 1024, min: 464 },
			items: 2,
		},
		mobile: {
			breakpoint: { max: 464, min: 0 },
			items: 1,
		},
	};

	const getTagComponents = () => {
		return tags.map((tag) => {
			return (
				<TagThumbnail
					key={tag.id}
					tag={tag}
					onTagClicked={(tag) => {
						if (duringDrag) {
							return;
						}

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
				<Box height={theme.iconSize(3)}>
					<Carousel
						responsive={responsive}
						infinite={true}
						slidesToSlide={10}
						arrows={false}
						centerMode={true}
						afterChange={(e) => {
							console.log('After change ' + e);
							setDuringDrag(false);
						}}
						beforeChange={(e) => {
							console.log('Before ' + e);
							setDuringDrag(true);
						}}
					>
						{getTagComponents()}
					</Carousel>
				</Box>
			)}
		</>
	);
}

export default TagThumbnails;
