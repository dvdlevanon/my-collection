import { Stack, ToggleButton, ToggleButtonGroup } from '@mui/material';
import TagChips from './TagChips';

function GalleryFilters({ conditionType, activeTags, onTagClick, onTagDelete, onChangeCondition }) {
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
						return tag.selected;
					}}
				/>
			)}
		</Stack>
	);
}

export default GalleryFilters;
