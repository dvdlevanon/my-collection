import { Stack } from '@mui/material';
import TagThumbnail from './TagThumbnail';

function TagThumbnails({ tags, onTagRemoved, onEditThumbnail }) {
	return (
		<Stack flexDirection="row" gap="10px" alignItems="center" flexWrap>
			{tags.map((tag) => {
				return (
					<TagThumbnail
						key={tag.id}
						tag={tag}
						onTagRemoved={onTagRemoved}
						onEditThumbnail={onEditThumbnail}
					/>
				);
			})}
		</Stack>
	);
}

export default TagThumbnails;
