import { Box } from '@mui/material';
import React from 'react';
import TagAnnotation from './TagAnnotation';

function TagsAnnotations({ annotations, selectedAnnotations, setSelectedAnnotations }) {
	const isSelectedAnnotation = (annotation) => {
		return selectedAnnotations.some((cur) => annotation.id == cur.id);
	};

	const annotationSelected = (e, annotation) => {
		if (isSelectedAnnotation(annotation)) {
			setSelectedAnnotations(selectedAnnotations.filter((cur) => annotation.id != cur.id));
		} else {
			setSelectedAnnotations([...selectedAnnotations, annotation]);
		}
	};

	return (
		<Box
			sx={{
				display: 'flex',
				flexDirection: 'row',
				flexWrap: 'wrap',
			}}
		>
			{annotations
				.sort((a, b) => {
					if (a.title == 'None' || b.title == 'None') {
						return 2;
					}

					if (a.title > b.title) {
						return 1;
					} else if (a.title < b.title) {
						return -1;
					} else {
						return 0;
					}
				})
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
	);
}

export default TagsAnnotations;
