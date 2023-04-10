import { Stack, ToggleButton, ToggleButtonGroup } from '@mui/material';
import TagChips from '../tags-chip/TagChips';

function GalleryFilters({ conditionType, activeTags, selectedTags, onTagClick, onTagDelete, onChangeCondition }) {
	return (
		<Stack flexDirection="row" gap="10px">
			{activeTags.length > 1 && (
				<Stack justifyContent="center" alignContent="center">
					<ToggleButtonGroup
						size="small"
						exclusive
						value={conditionType}
						onChange={(e, newValue) => {
							onChangeCondition(newValue);
						}}
					>
						<ToggleButton value="||">OR</ToggleButton>
						<ToggleButton value="&&">ADD</ToggleButton>
					</ToggleButtonGroup>
				</Stack>
			)}
			{activeTags.length > 0 && (
				<TagChips
					tags={activeTags}
					onClick={onTagClick}
					onDelete={onTagDelete}
					tagHighlightedPredicate={(tag) => {
						return selectedTags.some((cur) => cur.id == tag.id);
					}}
				/>
			)}
		</Stack>
	);
}

export default GalleryFilters;
