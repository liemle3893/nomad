import { inject as service } from '@ember/service';
import { computed } from '@ember/object';
import { alias, readOnly } from '@ember/object/computed';
import Controller, { inject as controller } from '@ember/controller';
import SortableFactory from 'nomad-ui/mixins/sortable-factory';

export default Controller.extend(SortableFactory([]), {
  userSettings: service(),
  pluginsController: controller('csi/plugins'),

  isForbidden: alias('pluginsController.isForbidden'),

  queryParams: {
    currentPage: 'page',
    sortProperty: 'sort',
    sortDescending: 'desc',
  },

  currentPage: 1,
  pageSize: readOnly('userSettings.pageSize'),

  sortProperty: 'id',
  sortDescending: false,

  listToSort: alias('model'),
  sortedPlugins: alias('listSorted'),

  // TODO: Remove once this page gets search capability
  resetPagination() {
    if (this.currentPage != null) {
      this.set('currentPage', 1);
    }
  },

  actions: {
    gotoPlugin(plugin) {
      this.transitionToRoute('csi.plugins.plugin', plugin.plainId);
    },
  },
});
