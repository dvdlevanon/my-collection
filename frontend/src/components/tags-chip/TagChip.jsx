import { Chip, Link } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import GalleryUrlParams from '../../utils/gallery-url-params';

function TagChip({ tag, linkable, onClick, onDelete, tagHighlightedPredicate }) {
	const tagChip = () => {
		return (
			<Chip
				color={tagHighlightedPredicate(tag) ? 'primary' : 'default'}
				label={tag.title}
				onClick={(e) => {
					e.stopPropagation();
					onClick(tag);
				}}
				onDelete={(e) => {
					e.preventDefault();
					e.stopPropagation();
					onDelete(tag);
				}}
			/>
		);
	};

	return (
		<>
			{(linkable && (
				<Link component={RouterLink} to={'/?' + GalleryUrlParams.buildUrlParams(tag.id)}>
					{tagChip()}
				</Link>
			)) ||
				tagChip()}
		</>
	);
}

export default TagChip;
