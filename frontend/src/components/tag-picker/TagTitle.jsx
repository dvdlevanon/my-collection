import { useTheme } from '@emotion/react';
import { Typography } from '@mui/material';
import TagsUtil from '../../utils/tags-util';

function TagTitle({ tag }) {
	const theme = useTheme();

	return (
		<Typography
			sx={{
				padding: theme.spacing(1),
				textAlign: 'center',
			}}
			noWrap
			variant="caption"
			textAlign={'start'}
		>
			{tag.title + ' (' + TagsUtil.itemsCount(tag) + ')'}
		</Typography>
	);
}

export default TagTitle;
