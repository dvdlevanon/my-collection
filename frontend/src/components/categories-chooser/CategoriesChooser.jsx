import { FormControl, MenuItem, Select } from '@mui/material';
import React, { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';

function CategoriesChooser({
	selectedIds,
	setCategories,
	onCreateCategoryClicked,
	allowToCreate,
	multiselect,
	placeholder,
}) {
	const unknownCategoryId = -1;
	const createCategoryId = -2;
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	// const [selectedCategories, setSelectedCategories] = useState(selectedIds);
	const [open, setOpen] = useState(false);

	const getUnknownCategory = () => {
		return {
			id: unknownCategoryId,
			title: placeholder,
		};
	};

	const getCreateCategory = () => {
		return {
			id: createCategoryId,
			title: 'Create New...',
		};
	};

	const getCategories = () => {
		let result = [];
		result.push(getUnknownCategory());
		if (allowToCreate) {
			result.push(getCreateCategory());
		}
		result = result.concat(TagsUtil.getCategories(tagsQuery.data));
		result = result.filter((category) => !TagsUtil.isSpecialCategory(category.id));
		return result;
	};

	const onChange = (event) => {
		let value = event.target.value;

		if (value[value.length - 1] == createCategoryId) {
			onCreateCategoryClicked();
			return;
		}

		if (value[value.length - 1] == unknownCategoryId) {
			setCategories([]);
			// setSelectedCategories([unknownCategoryId]);
			setOpen(false);
			return;
		}

		if (value.length > 1 && value.some((category) => category == unknownCategoryId)) {
			value = value.filter((category) => category != unknownCategoryId);
		}

		setCategories(value);
		// setSelectedCategories(value);
		setOpen(false);
	};

	return (
		<FormControl fullWidth>
			<Select
				open={open}
				onOpen={() => setOpen(true)}
				onClose={() => setOpen(false)}
				onChange={onChange}
				size="small"
				multiple={multiselect}
				value={selectedIds}
				displayEmpty
			>
				{getCategories().map((category) => {
					return (
						<MenuItem key={category.id} value={category.id}>
							{category.title}
						</MenuItem>
					);
				})}
			</Select>
		</FormControl>
	);
}

export default CategoriesChooser;
