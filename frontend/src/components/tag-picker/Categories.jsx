import { Stack } from '@mui/material';
import React from 'react';
import Category from './Category';

function Categories({ categories, onCategoryClicked }) {
	return (
		<Stack flexDirection="row">
			{categories.map((category) => {
				return <Category key={category.id} category={category} onClick={onCategoryClicked} />;
			})}
		</Stack>
	);
}

export default Categories;
