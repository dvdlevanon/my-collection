import { useMutation, useQueryClient } from '@tanstack/react-query';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import { useTagAnnotationsStore } from './tagAnnotationsStore';

export const useAnnotations = () => {
	const queryClient = useQueryClient();
	const tagAnnotationsStore = useTagAnnotationsStore();
	const addAnnotationToTagMutation = useMutation({ mutationFn: Client.addAnnotationToTag });
	const removeAnnotationFromTagMutation = useMutation({ mutationFn: Client.removeAnnotationFromTag });

	const addAnnotation = (annotation, successCallback) => {
		let payload = { tagId: tagAnnotationsStore.tag.id, annotation: annotation };
		addAnnotationToTagMutation.mutate(payload, {
			onSuccess: (response) => {
				response.json().then((tagAnnotation) => {
					ReactQueryUtil.updateTags(queryClient, tagAnnotationsStore.tag.id, (draft) => {
						if (!draft.tags_annotations) {
							draft.tags_annotations = [];
						}
						draft.tags_annotations.push(tagAnnotation);
					});

					postUpdate();
					if (successCallback) successCallback();
				});
			},
		});
	};

	const removeAnnotation = (annotation) => {
		let payload = { tagId: tagAnnotationsStore.tag.id, annotationId: annotation.id };
		removeAnnotationFromTagMutation.mutate(payload, {
			onSuccess: (response, params) => {
				ReactQueryUtil.updateTags(queryClient, tagAnnotationsStore.tag.id, (draft) => {
					draft.tags_annotations = draft.tags_annotations.filter((annotation) => {
						return annotation.id != params.annotationId;
					});
				});

				postUpdate();
			},
		});
	};

	const postUpdate = () => {
		queryClient.invalidateQueries(ReactQueryUtil.availableAnnotationsQuery(tagAnnotationsStore.tag.parentId));
		const updatedTag = ReactQueryUtil.getTag(queryClient, tagAnnotationsStore.tag.id);
		if (updatedTag) {
			tagAnnotationsStore.setTag(updatedTag);
		}
	};

	return { addAnnotation, removeAnnotation };
};
