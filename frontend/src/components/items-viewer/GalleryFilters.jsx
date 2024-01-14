import { useTheme } from '@emotion/react';
import CloseIcon from '@mui/icons-material/Close';
import { IconButton, Paper, Stack, ToggleButton, ToggleButtonGroup, Tooltip } from '@mui/material';
import TagChips from '../tags-chip/TagChips';
import TextFieldWithKeyboard from '../text-field-with-keyboard/TextFieldWithKeyboard';

function GalleryFilters({
	activeTags,
	selectedTags,
	onTagClick,
	onTagDelete,
	conditionType,
	setConditionType,
	searchTerm,
	setSearchTerm,
	galleryUrlParams,
}) {
	const theme = useTheme();

	return (
		<Stack flexDirection="row" gap={theme.spacing(1)} alignItems="center">
			<TextFieldWithKeyboard
				variant="outlined"
				autoFocus
				label="Search..."
				type="search"
				size="small"
				onChange={(value) => setSearchTerm(value)}
				value={searchTerm}
			></TextFieldWithKeyboard>
			{activeTags.length > 0 && (
				<Paper
					variant="elevation"
					elevation={4}
					sx={{
						display: 'flex',
						gap: theme.spacing(1),
						padding: theme.spacing(0),
					}}
				>
					<Tooltip title="Remove all filters">
						<IconButton
							onClick={() => {
								galleryUrlParams.deactivateAllTags();
							}}
						>
							<CloseIcon sx={{ fontSize: theme.iconSize(1) }} />
						</IconButton>
					</Tooltip>
					<TagChips
						tags={activeTags}
						linkable={false}
						onClick={onTagClick}
						onDelete={onTagDelete}
						tagHighlightedPredicate={(tag) => {
							return selectedTags.some((cur) => cur.id == tag.id);
						}}
					/>
					<Stack justifyContent="center" alignContent="center">
						<ToggleButtonGroup
							size="small"
							exclusive
							value={conditionType}
							onChange={(e, newValue) => {
								setConditionType(newValue);
							}}
						>
							<ToggleButton value="||">OR</ToggleButton>
							<ToggleButton value="&&">AND</ToggleButton>
						</ToggleButtonGroup>
					</Stack>
				</Paper>
			)}
		</Stack>
	);
}

export default GalleryFilters;
