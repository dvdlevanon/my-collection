import { Stack, ToggleButton, ToggleButtonGroup } from '@mui/material';
import TagChips from '../tags-chip/TagChips';

function GalleryFilters({
	itemsSizeName,
	setItemsSizeName,
	conditionType,
	activeTags,
	onTagClick,
	onTagDelete,
	onChangeCondition,
}) {
	return (
		<Stack flexDirection="row" gap="10px">
			<Stack justifyContent="center" alignContent="center">
				<ToggleButtonGroup
					size="small"
					exclusive
					value={itemsSizeName}
					onChange={(e, newValue) => {
						setItemsSizeName(newValue);
					}}
				>
					<ToggleButton value="xs">xs</ToggleButton>
					<ToggleButton value="s">s</ToggleButton>
					<ToggleButton value="m">m</ToggleButton>
					<ToggleButton value="l">l</ToggleButton>
					<ToggleButton value="xl">xl</ToggleButton>
				</ToggleButtonGroup>
			</Stack>
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
