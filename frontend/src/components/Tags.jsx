import { TextField } from '@mui/material';
import { Box } from '@mui/system';
import { useState } from 'react';
import { useQuery } from 'react-query';
import Client from '../network/client';
import ReactQueryUtil from '../utils/react-query-util';
import Tag from './Tag';
import TagAnnotation from './TagAnnotation';

function Tags({ tags, parentId, size, onTagSelected }) {
	let [searchTerm, setSearchTerm] = useState('');
	let [selectedAnnotaions, setSelectedAnnotations] = useState([]);
	const availableAnnotationsQuery = useQuery({
		queryKey: ReactQueryUtil.availableAnnotationsKey(parentId),
		queryFn: () => Client.getAvailableAnnotations(parentId),
		onSuccess: (availableAnnotations) => {
			setSelectedAnnotations(
				selectedAnnotaions.filter((selected) => {
					return availableAnnotations.some((annotation) => selected.id == annotation.id);
				})
			);
		},
	});

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
		<Box
			sx={{
				position: 'absolute',
				zIndex: '100',
				top: '0px',
				left: '0px',
				right: '0px',
			}}
		>
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
					padding: '10px',
				}}
			>
				<TextField
					variant="outlined"
					autoFocus
					label="Search..."
					type="search"
					size="small"
					onChange={(e) => onSearchTermChanged(e)}
				></TextField>
				<Box
					sx={{
						display: 'flex',
						flexDirection: 'row',
					}}
				>
					{availableAnnotationsQuery.isSuccess &&
						availableAnnotationsQuery.data
							.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0))
							.map((annotation) => {
								return (
									<TagAnnotation
										key={annotation.id}
										selectedAnnotaions
										annotation={annotation}
										selected={isSelectedAnnotation(annotation)}
										onClick={annotationSelected}
									/>
								);
							})}
				</Box>
			</Box>
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
					padding: '10px',
					gap: '10px',
					flexWrap: 'wrap',
				}}
			>
				{filterTags().map((tag) => {
					return (
						<div key={tag.id}>
							<Tag tag={tag} size={size} onTagSelected={onTagSelected} />
						</div>
					);
				})}
			</Box>
		</Box>
	);
}

export default Tags;
