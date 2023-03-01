import { FormControl, MenuItem, Select } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';

function DirectoryCategoriesCell({ directory, setCategories, onCreateCategoryClicked }) {
	const unknownCategoryId = -1;
	const createCategoryId = -2;
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const [selectedCategories, setSelectedCategories] = useState([unknownCategoryId]);
	const [open, setOpen] = React.useState(false);

	useEffect(() => {
		if (directory.tags != null) {
			setSelectedCategories(directory.tags.map((dir) => dir.id));
		}
	}, []);

	const getUnknownCategory = () => {
		return {
			id: unknownCategoryId,
			title: 'Not Specified',
		};
	};

	const getCreateCategory = () => {
		return {
			id: createCategoryId,
			title: 'Create New...',
		};
	};

	const getCategories = () => {
		let result = TagsUtil.getCategories(tagsQuery.data);
		result = result.filter((category) => !TagsUtil.isDirectoriesCategory(category.id));
		result.push(getUnknownCategory());
		result.push(getCreateCategory());
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
			setSelectedCategories([unknownCategoryId]);
			setOpen(false);
			return;
		}

		if (value.length > 1 && value.some((category) => category == unknownCategoryId)) {
			value = value.filter((category) => category != unknownCategoryId);
		}

		setCategories(value);
		setSelectedCategories(value);
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
				multiple
				value={selectedCategories}
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

export default DirectoryCategoriesCell;
