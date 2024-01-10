import { Stack } from '@mui/material';
import TagThumbnail from './TagThumbnail';

function TagThumbnails({ tags, onTagClicked, onTagRemoved, onEditThumbnail, withRemoveOption, additionSx }) {
	return (
		<Stack flexDirection="row" gap="10px" alignItems="center" sx={additionSx}>
			{tags.map((tag) => {
				return (
					<TagThumbnail
						key={tag.id}
						tag={tag}
						onTagClicked={onTagClicked}
						onTagRemoved={onTagRemoved}
						onEditThumbnail={onEditThumbnail}
						withRemoveOption={withRemoveOption}
					/>
				);
			})}
		</Stack>
	);
}

export default TagThumbnails;
