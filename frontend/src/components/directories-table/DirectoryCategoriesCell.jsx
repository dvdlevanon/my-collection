import React, { useEffect, useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import CategoriesChooser from '../categories-chooser/CategoriesChooser';

function DirectoryCategoriesCell({ directory, setCategories, onCreateCategoryClicked }) {
	const unknownCategoryId = -1;
	const tagsQuery = useQuery(ReactQueryUtil.TAGS_KEY, Client.getTags);
	const [selectedCategories, setSelectedCategories] = useState([unknownCategoryId]);

	useEffect(() => {
		if (directory.tags != null) {
			setSelectedCategories(directory.tags.map((dir) => dir.id));
		}
	}, []);

	return (
		<CategoriesChooser
			multiselect={true}
			allowToCreate={true}
			placeholder="Not Specified"
			selectedIds={selectedCategories}
			setCategories={(categoriyIds) => {
				setCategories(categoriyIds);
				setSelectedCategories(categoriyIds);
			}}
			onCreateCategoryClicked={onCreateCategoryClicked}
		/>
	);
}

export default DirectoryCategoriesCell;
