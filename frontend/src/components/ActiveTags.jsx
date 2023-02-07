import { Stack } from '@mui/material';
import ActiveTag from './ActiveTag';

function ActiveTags({ activeTags, onTagDeactivated, onTagSelected, onTagDeselected }) {
	return (
		<Stack flexDirection="row" gap="10px" flexWrap>
			{activeTags.map((tag) => {
				return (
					<ActiveTag
						key={tag.id}
						tag={tag}
						onTagDeactivated={onTagDeactivated}
						onTagSelected={onTagSelected}
						onTagDeselected={onTagDeselected}
					/>
				);
			})}
		</Stack>
	);
}

export default ActiveTags;
