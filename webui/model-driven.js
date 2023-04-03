// load the schema (YANG) and the object for editing
function load(module) {
    const parent = document.getElementById("form");
    Promise.allSettled([
        fetch(`/restconf/data/${module}:`)
            .then(resp => resp.json()),
        fetch(`/restconf/schema/${module}/`, {
            method: 'GET',
            headers: {'Accept': 'application/json'}
        }).then(resp => resp.json()),
    ]).then((responses) => {
        render(parent, responses[0].value, responses[1].value.module);
    });
}

// navigate thru the meta along with the object to build forms but
// as a pattern for all model driven UI
function render(parent, obj, meta) {
    const editableFields = meta.dataDef.filter(def => {
        // ignore mertics
        if (def.leaf?.config == false) {
            return false;
        }
        // would normally recurse here
        if ('list' in def || 'container' in def) { 
            return false;
        }
        return true;
    })
    // here you would normally adjust the input type based on details
    // in the 'def' object like number v.s. string, etc.
    parent.innerHTML = `
        <table>
        ${editableFields.map(def => `
            <tr>
                <td><label>${def.ident}</label></td>
                <td><input value="${obj[def.ident] || ''}"></td>
            </tr>
        </table>`).join('')}
    `;

    // recurse into lists and containers here
}

// car is not a very exciting object to edit but still demonstrates feature
load('car');
