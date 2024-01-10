import { useTheme } from '@emotion/react';
import { Stack } from '@mui/material';
import TagChip from './TagChip';

function TagChips({ tags, linkable, onClick, onDelete, tagHighlightedPredicate }) {
	const theme = useTheme();

	return (
		<Stack flexDirection="row" gap={theme.spacing(1)} alignItems="center" flexWrap>
			{tags.map((tag) => {
				return (
					<TagChip
						key={tag.id}
						tag={tag}
						linkable={linkable}
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
