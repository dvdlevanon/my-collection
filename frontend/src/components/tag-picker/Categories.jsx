import { useTheme } from '@emotion/react';
import AddIcon from '@mui/icons-material/Add';
import { IconButton, Stack, Tooltip } from '@mui/material';
import React, { useState } from 'react';
import TagsUtil from '../../utils/tags-util';
import AddTagDialog from '../dialogs/AddTagDialog';
import Category from './Category';

function Categories({ categories, onCategoryClicked, selectedCategoryId }) {
	const [addCategoryDialogOpened, setAddCategoryDialogOpened] = useState(false);
	const theme = useTheme();

	const sortedCategories = () => {
		return categories.sort((cat1, cat2) => {
			return TagsUtil.isSpecialCategory(cat1.id) ? -1 : 1;
		});
	};

	return (
		<Stack padding={theme.spacing(1)} flexDirection="row" alignItems="center" gap={theme.spacing(1)}>
			<Tooltip title="Add Category">
				<IconButton size="small" onClick={() => setAddCategoryDialogOpened(true)}>
					<AddIcon sx={{ fontSize: theme.iconSize(1) }} />
				</IconButton>
			</Tooltip>
			{sortedCategories().map((category) => {
				return (
					<Category
						isHighlighted={category.id == selectedCategoryId}
						key={category.id}
						category={category}
						onClick={onCategoryClicked}
					/>
				);
			})}
			{categories.length == 0 && <div>Add a category</div>}
			<AddTagDialog
				open={addCategoryDialogOpened}
				parentId={null}
				verb="Category"
				onClose={() => setAddCategoryDialogOpened(false)}
			/>
		</Stack>
	);
}

export default Categories;
