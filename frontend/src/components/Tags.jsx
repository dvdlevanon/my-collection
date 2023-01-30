import { TextField } from '@mui/material';
import { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import Tag from './Tag';
import TagAnnotation from './TagAnnotation';
import styles from './Tags.module.css';

function Tags({ tags, parentId, size, onTagSelected }) {
	let [searchTerm, setSearchTerm] = useState('');
	let [selectedAnnotaions, setSelectedAnnotations] = useState([]);
	const availableAnnotationsQuery = useQuery(ReactQueryUtil.availableAnnotationsKey(parentId), () =>
		Client.getAvailableAnnotations(parentId)
	);

	const onSearchTermChanged = (e) => {
		setSearchTerm(e.target.value);
	};

	const filterTagsBySearch = (tags) => {
		let filteredTags = tags;

		if (searchTerm) {
			filteredTags = tags.filter((tag) => {
				return tag.title.toLowerCase().includes(searchTerm.toLowerCase());
			});
		}

		return filteredTags;
	};

	const filterTagsByAnnotations = (tags) => {
		return tags.filter((cur) => {
			if (selectedAnnotaions.length == 0) {
				return true;
			}

			if (!cur.tags_annotations) {
				return false;
			}

			return cur.tags_annotations.some((tagAnnotation) => {
				return selectedAnnotaions.some((annotation) => annotation.id == tagAnnotation.id);
			});
		});
	};

	const filterTags = () => {
		let filteredTags = filterTagsByAnnotations(filterTagsBySearch(tags));

		return filteredTags.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0));
		// let shuffeledTags = filteredTags.sort(() => (Math.random() > 0.5 ? 1 : -1));
		// return shuffeledTags;
	};

	const isSelectedAnnotation = (annotation) => {
		return selectedAnnotaions.some((cur) => annotation.id == cur.id);
	};

	const annotationSelected = (e, annotation) => {
		if (isSelectedAnnotation(annotation)) {
			setSelectedAnnotations(selectedAnnotaions.filter((cur) => annotation.id != cur.id));
		} else {
			setSelectedAnnotations([...selectedAnnotaions, annotation]);
		}
	};

	return (
		<div className={styles.tags_holder}>
			<div className={styles.filters_holder}>
				<TextField
					variant="outlined"
					autoFocus
					fullWidth
					label="Search..."
					type="search"
					sx={{ width: '500px' }}
					onChange={(e) => onSearchTermChanged(e)}
				></TextField>
				{availableAnnotationsQuery.isSuccess &&
					availableAnnotationsQuery.data
						.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0))
						.map((annotation) => {
							return (
								<TagAnnotation
									key={annotation.id}
									annotation={annotation}
									selected={isSelectedAnnotation(annotation)}
									onClick={annotationSelected}
								/>
							);
						})}
			</div>
			<div className={styles.tags}>
				{filterTags().map((tag) => {
					return (
						<div key={tag.id}>
							<Tag tag={tag} size={size} onTagSelected={onTagSelected} />
						</div>
					);
				})}
			</div>
		</div>
	);
}

export default Tags;
