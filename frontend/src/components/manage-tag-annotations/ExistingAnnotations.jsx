import { Stack } from '@mui/material';
import TagAnnotation from '../tag-annotation/TagAnnotation';
import { useTagAnnotationsStore } from './tagAnnotationsStore';
import { useAnnotations } from './useAnnotations';

function ExistingAnnotations({ resetNewAnnotationText }) {
	const tagAnnotationsStore = useTagAnnotationsStore();
	const { addAnnotation, removeAnnotation } = useAnnotations();

	const addAnnotationClicked = (e, annotation) => {
		e.preventDefault();
		e.stopPropagation();
		addAnnotation(annotation, resetNewAnnotationText);
	};

	const removeAnnotationClicked = (e, annotation) => {
		e.preventDefault();
		e.stopPropagation();
		removeAnnotation(annotation);
	};

	return (
		<Stack flexDirection="row" flexWrap="wrap">
			{tagAnnotationsStore.getAttachedAnnotations().map((annotation) => {
				return (
					<TagAnnotation
						key={annotation.id}
						annotation={annotation}
						selected={true}
						onRemoveClicked={removeAnnotationClicked}
					/>
				);
			})}
			{tagAnnotationsStore.getUnattachedAnnotations().map((annotation) => {
				return (
					<TagAnnotation
						key={annotation.id}
						annotation={annotation}
						selected={false}
						onClick={addAnnotationClicked}
					/>
				);
			})}
		</Stack>
	);
}

export default ExistingAnnotations;
