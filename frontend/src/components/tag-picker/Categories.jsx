import AddIcon from '@mui/icons-material/Add';
import { IconButton, Stack, Tooltip } from '@mui/material';
import React, { useState } from 'react';
import AddTagDialog from '../dialogs/AddTagDialog';
import Category from './Category';

function Categories({ categories, onCategoryClicked, selectedCategoryId }) {
	const [addCategoryDialogOpened, setAddCategoryDialogOpened] = useState(false);

	return (
		<Stack flexDirection="row" alignItems="center" gap="10px">
			<Tooltip title="Add Category">
				<IconButton size="small" onClick={() => setAddCategoryDialogOpened(true)}>
					<AddIcon sx={{ fontSize: '20px' }} />
				</IconButton>
			</Tooltip>
			{categories.map((category) => {
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
			{addCategoryDialogOpened && (
				<AddTagDialog parentId={null} verb="Category" onClose={() => setAddCategoryDialogOpened(false)} />
			)}
		</Stack>
	);
}

export default Categories;
