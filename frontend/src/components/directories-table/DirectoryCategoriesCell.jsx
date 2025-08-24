import { useQuery } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import CategoriesChooser from '../categories-chooser/CategoriesChooser';

function DirectoryCategoriesCell({ directory, setCategories }) {
	const unknownCategoryId = -1;
	const tagsQuery = useQuery({ queryKey: ReactQueryUtil.TAGS_KEY, queryFn: Client.getTags });
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
		/>
	);
}

export default DirectoryCategoriesCell;
