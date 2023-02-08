import AddLink from '@mui/icons-material/AddLink';
import { default as RemoveIcon } from '@mui/icons-material/Close';
import ImageIcon from '@mui/icons-material/Image';
import OptionsIcon from '@mui/icons-material/Tune';
import { SpeedDial, SpeedDialAction } from '@mui/material';
import React from 'react';

function TagSpeedDial({ tag, hidden, onManageAttributesClicked, onRemoveTagClicked, onManageImageClicked }) {
	return (
		<React.Fragment>
			{!hidden && (
				<SpeedDial
					sx={{
						position: 'absolute',
						bottom: '0px',
						right: '0px',
						padding: '5px',
						'& .MuiFab-primary': {
							width: 40,
							height: 40,
							backgroundColor: 'primary.main',
						},
					}}
					ariaLabel="tag-actions"
					icon={<OptionsIcon />}
					onClick={(e) => e.stopPropagation()}
				>
					<SpeedDialAction
						key="manage-image"
						tooltipTitle="Image options"
						icon={<ImageIcon />}
						onClick={(e) => {
							onManageImageClicked();
						}}
					/>

					<SpeedDialAction
						key="manage-annotations"
						tooltipTitle="Manage annotations"
						icon={<AddLink />}
						onClick={(e) => {
							onManageAttributesClicked(e);
						}}
					/>
					<SpeedDialAction
						key="remove-tag"
						tooltipTitle="Remove tag"
						icon={<RemoveIcon />}
						onClick={(e) => {
							e.stopPropagation();
							onRemoveTagClicked();
						}}
					/>
				</SpeedDial>
			)}
			{/* <Box
				component="input"
				onClick={(e) => e.stopPropagation()}
				onChange={(e) => imageSelected(e)}
				accept="image/*"
				id="choose-file"
				type="file"
				sx={{
					display: 'none',
				}}
				ref={fileDialog}
			/> */}
		</React.Fragment>
	);
}
// let [removeTagDialogOpened, setRemoveTagDialogOpened] = useState(false);

export default TagSpeedDial;
