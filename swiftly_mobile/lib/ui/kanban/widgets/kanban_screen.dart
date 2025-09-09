import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../utils/responsive_layout.dart';
import 'mobile/kanban_screen_mobile.dart';
import 'desktop/kanban_screen_desktop.dart';

class KanbanScreen extends ConsumerWidget {
  const KanbanScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return const ResponsiveLayout(
      mobile: KanbanScreenMobile(),
      desktop: KanbanScreenDesktop(),
    );
  }
}
