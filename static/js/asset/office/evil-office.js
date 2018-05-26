import { LitElement, html } from '../lit-element.bundle.js';

class EvilOffice extends LitElement {

    static get properties() {
        return {
            mood: String
        }
    }

    _render({ mood }) {
        return html `<style> .mood { color: green; } </style>
      Web Components are <span class="mood">${mood}</span>!`;
    }

}

customElements.define('evil-office', EvilOffice);