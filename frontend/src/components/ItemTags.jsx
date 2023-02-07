import AddIcon from '@mui/icons-material/Add';
import { Box, Fab, Stack } from '@mui/material';
import ItemTag from './ItemTag';

function ItemTags({ item, onAddTag, onTagRemoved }) {
	return (
		<Stack>
			<Stack gap="10px" padding="10px">
				{item.tags &&
					item.tags.map((tag) => {
						return <ItemTag key={tag.id} tag={tag} onRemoveClicked={onTagRemoved} />;
					})}
			</Stack>
			<Box
				sx={{
					padding: '10px',
					position: 'absolute',
					bottom: '0px',
					right: '0px',
				}}
			>
				<Fab
					onClick={(e) => {
						onAddTag();
					}}
				>
					<AddIcon />
				</Fab>
			</Box>
		</Stack>
	);
}

export default ItemTags;
