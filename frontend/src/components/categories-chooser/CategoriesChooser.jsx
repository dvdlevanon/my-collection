import { FormControl, MenuItem, Select, Stack } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useState } from 'react';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import AddTagDialog from '../dialogs/AddTagDialog';

function CategoriesChooser({ selectedIds, setCategories, allowToCreate, multiselect, placeholder }) {
	const unknownCategoryId = -1;
	const createCategoryId = -2;
	const [open, setOpen] = useState(false);
	const [addCategoryDialogOpened, setAddCategoryDialogOpened] = useState(false);
	const tagsQuery = useQuery({ queryKey: ReactQueryUtil.TAGS_KEY, queryFn: Client.getTags });

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
			setAddCategoryDialogOpened(true);
			return;
		}

		if (value[value.length - 1] == unknownCategoryId) {
			setCategories([]);
			setOpen(false);
			return;
		}

		if (value.length > 1 && value.some((category) => category == unknownCategoryId)) {
			value = value.filter((category) => category != unknownCategoryId);
		}

		setCategories(value);
		setOpen(false);
	};

	return (
		<Stack>
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
			{addCategoryDialogOpened && (
				<AddTagDialog
					open={addCategoryDialogOpened}
					parentId={null}
					verb="Category"
					onClose={() => setAddCategoryDialogOpened(false)}
				/>
			)}
		</Stack>
	);
}

export default CategoriesChooser;
