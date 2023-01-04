
function SuperTag({ superTag, onSuperTagSelected, onSuperTagDeselected }) {

    const superTagClicked = (e) => {
        if (superTag.selected) {
            onSuperTagDeselected(superTag)
        } else {
            onSuperTagSelected(superTag)
        }
    }

    return (
        <div className="super-tag" onClick={(e) => superTagClicked(e)}>
            {superTag.title}
        </div>
    )
}

export default SuperTag