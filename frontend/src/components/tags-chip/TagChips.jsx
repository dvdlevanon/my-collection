import { Stack } from '@mui/material';
import TagChip from './TagChip';

function TagChips({ tags, onClick, onDelete, tagHighlightedPredicate }) {
	return (
		<Stack flexDirection="row" gap="10px" alignItems="center" flexWrap>
			{tags.map((tag) => {
				return (
					<TagChip
						key={tag.id}
						tag={tag}
						onClick={onClick}
						onDelete={onDelete}
						tagHighlightedPredicate={tagHighlightedPredicate}
					/>
				);
			})}
		</Stack>
	);
}

export default TagChips;
