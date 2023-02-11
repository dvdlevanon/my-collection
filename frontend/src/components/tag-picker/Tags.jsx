import AddIcon from '@mui/icons-material/Add';
import { Button, IconButton, Link, Stack, TextField } from '@mui/material';
import { Box } from '@mui/system';
import { useState } from 'react';
import { useQuery } from 'react-query';
import { Link as RouterLink } from 'react-router-dom';
import Client from '../../network/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import AddTagDialog from '../dialogs/AddTagDialog';
import Tag from './Tag';
import TagAnnotation from './TagAnnotation';

function Tags({ tags, parentId, size, onTagSelected }) {
	const [addTagDialogOpened, setAddTagDialogOpened] = useState(false);
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
		<Stack width={'100%'} height={'100%'} flexGrow={1}>
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
					padding: '10px',
					gap: '10px',
				}}
			>
				{!TagsUtil.isDirectoriesCategory(parentId) && (
					<IconButton onClick={() => setAddTagDialogOpened(true)}>
						<AddIcon />
					</IconButton>
				)}
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
					flexGrow: 1,
				}}
			>
				{filterTags().map((tag) => {
					return <Tag key={tag.id} tag={tag} size={size} onTagSelected={onTagSelected} />;
				})}
				{!TagsUtil.isDirectoriesCategory(parentId) && (
					<Tag key="add-tag" tag={{ id: -1 }} size={size} onTagSelected={() => setAddTagDialogOpened(true)} />
				)}

				{TagsUtil.isDirectoriesCategory(parentId) && tags.length == 0 && (
					<Stack
						direction="row"
						gap="10px"
						justifyContent="center"
						alignItems="center"
						sx={{
							width: '100%',
							height: '100%',
						}}
					>
						No diretories found
						<Link component={RouterLink} to="spa/manage-directories">
							<Button variant="outlined">Manage Directories</Button>
						</Link>
					</Stack>
				)}
			</Box>
			{addTagDialogOpened && (
				<AddTagDialog parentId={parentId} verb="Tag" onClose={() => setAddTagDialogOpened(false)} />
			)}
		</Stack>
	);
}

export default Tags;
