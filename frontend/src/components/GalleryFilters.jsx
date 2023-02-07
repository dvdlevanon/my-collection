import { Stack, ToggleButton, ToggleButtonGroup } from '@mui/material';
import ActiveTags from './ActiveTags';

function GalleryFilters({
	conditionType,
	activeTags,
	onTagDeactivated,
	onTagSelected,
	onTagDeselected,
	onChangeCondition,
}) {
	const onConditionChanged = (e, newValue) => {
		onChangeCondition(newValue);
	};

	return (
		<Stack flexDirection="row" gap="10px">
			{activeTags.length > 1 && (
				<Stack justifyContent="center" alignContent="center">
					<ToggleButtonGroup size="small" exclusive value={conditionType} onChange={onConditionChanged}>
						<ToggleButton value="||">OR</ToggleButton>
						<ToggleButton value="&&">ADD</ToggleButton>
					</ToggleButtonGroup>
				</Stack>
			)}
			{activeTags.length > 0 && (
				<ActiveTags
					activeTags={activeTags}
					onTagDeactivated={onTagDeactivated}
					onTagSelected={onTagSelected}
					onTagDeselected={onTagDeselected}
				/>
			)}
		</Stack>
	);
}

export default GalleryFilters;
