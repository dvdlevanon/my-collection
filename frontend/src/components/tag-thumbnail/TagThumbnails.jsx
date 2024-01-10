import { Box, Stack } from '@mui/material';
import { useState } from 'react';
import Carousel from 'react-multi-carousel';
import 'react-multi-carousel/lib/styles.css';
import TagThumbnail from './TagThumbnail';

function TagThumbnails({ tags, carouselMode, onTagClicked, onTagRemoved, onEditThumbnail, withRemoveOption }) {
	const [duringDrag, setDuringDrag] = useState(false);

	const responsive = {
		superLargeDesktop: {
			breakpoint: { max: 4000, min: 3000 },
			items: 30,
		},
		desktop: {
			breakpoint: { max: 3000, min: 1024 },
			items: 30,
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
					onEditThumbnail={onEditThumbnail}
					withRemoveOption={withRemoveOption}
					isLink={onTagClicked == null}
				/>
			);
		});
	};

	return (
		<>
			{!carouselMode && (
				<Stack flexDirection="row" gap="10px" alignItems="center">
					{getTagComponents()}
				</Stack>
			)}
			{carouselMode && (
				<Box height="100px">
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
