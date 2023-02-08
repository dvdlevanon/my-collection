import { Chip } from '@mui/material';

function TagChip({ tag, onClick, onDelete, tagHighlightedPredicate }) {
	return (
		<Chip
			color={tagHighlightedPredicate(tag) ? 'primary' : 'default'}
			label={tag.title}
			onClick={(e) => {
				e.stopPropagation();
				onClick(tag);
			}}
			onDelete={(e) => {
				e.stopPropagation();
				onDelete(tag);
			}}
		/>
	);
}

export default TagChip;
